package launcher

import (
	"context"
	"fmt"
	"net/http"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type spec[T any] struct {
	nsf  NewServiceFunc[T]
	wsf  WithStorageFunc[T]
	init ServiceInitFunc[T]
	opts []service.Option[T]
}

func (s spec[T]) newService(db persistence.Storage) (svc T, grpcOpts []server.StartGRPCServerOption, err error) {
	// Append the WithStorageFunc to the specified service options.
	var opts []service.Option[T]
	opts = append(opts, s.wsf(db))

	// Create the service with the NewServiceFunc using the supplied server options
	svc = s.nsf(opts...)

	// Initialize the service using the ServiceInitFunc. This returns a possible list of StartGRPCServerOptions that we need to return
	grpcOpts, err = s.init(svc)
	if err != nil {
		return *new(T), nil, err
	}

	return
}

func (s spec[T]) NewService(db persistence.Storage) (svc any, grpcOpts []server.StartGRPCServerOption, err error) {
	return s.newService(db)
}

type Launcher[T any] struct {
	srv       *server.Server
	component string
	db        persistence.Storage
	log       *logrus.Entry
	grpcOpts  []server.StartGRPCServerOption

	Service T
}

func (l Launcher[T]) ToAny() *Launcher[any] {
	return &Launcher[any]{
		srv:       l.srv,
		component: l.component,
		db:        l.db,
		log:       l.log,
		grpcOpts:  l.grpcOpts,
		Service:   l.Service,
	}
}

type NewServiceFunc[T any] func(opts ...service.Option[T]) T
type WithStorageFunc[T any] func(db persistence.Storage) service.Option[T]
type ServiceInitFunc[T any] func(svc T) ([]server.StartGRPCServerOption, error)

func NewLauncher[T any](component string, nsf NewServiceFunc[T], wsf WithStorageFunc[T], init ServiceInitFunc[T], serviceOpts ...service.Option[T]) (l *Launcher[T], err error) {
	l = new(Launcher[T])
	l.component = component

	// Print Clouditor header with the used Clouditor version
	printClouditorHeader(fmt.Sprintf("Clouditor %s Service", cases.Title(language.English).String(l.component)))

	// Set log level
	err = l.initLogging()
	if err != nil {
		return nil, err
	}

	// Set storage
	err = l.initStorage()
	if err != nil {
		return nil, err
	}

	// Build a service spec to be consistent with the multi-launcher approach
	spec := &spec[T]{
		nsf:  nsf,
		wsf:  wsf,
		init: init,
		opts: serviceOpts,
	}

	// Create a new service from the service spec
	l.Service, l.grpcOpts, err = spec.newService(l.db)
	if err != nil {
		return nil, err
	}

	return
}

// Launch starts the gRPC server and the corresponding gRPC-HTTP gateway with the given gRPC Server Options
func (l *Launcher[T]) Launch() (err error) {
	var (
		grpcPort uint16
		httpPort uint16
		restOpts []rest.ServerConfigOption
		grpcOpts []server.StartGRPCServerOption
	)

	// Default gRPC opts we want for all services
	grpcOpts = []server.StartGRPCServerOption{
		server.WithJWKS(viper.GetString(config.APIJWKSURLFlag)),
		server.WithReflection(),
	}

	// Append launch-specific ones
	grpcOpts = append(grpcOpts, l.grpcOpts...)

	grpcPort = viper.GetUint16(config.APIgRPCPortOrchestratorFlag)
	httpPort = viper.GetUint16(config.APIHTTPPortOrchestratorFlag)

	restOpts = []rest.ServerConfigOption{
		rest.WithAllowedOrigins(viper.GetStringSlice(config.APICORSAllowedOriginsFlags)),
		rest.WithAllowedHeaders(viper.GetStringSlice(config.APICORSAllowedHeadersFlags)),
		rest.WithAllowedMethods(viper.GetStringSlice(config.APICORSAllowedMethodsFlags)),
	}

	l.log.Infof("Starting gRPC endpoint on :%d", grpcPort)

	// Start the gRPC server
	_, l.srv, err = server.StartGRPCServer(
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		grpcOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to serve gRPC endpoint: %w", err)
	}

	// Start the gRPC-HTTP gateway
	err = rest.RunServer(context.Background(),
		grpcPort,
		httpPort,
		restOpts...,
	)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve gRPC-HTTP gateway: %v", err)
	}

	l.log.Infof("Stopping gRPC endpoint")
	l.srv.Stop()

	return
}

// initLogging initializes the logging
func (l *Launcher[T]) initLogging() error {
	l.log = logrus.WithField("component", l.component)
	l.log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}

	level, err := logrus.ParseLevel(viper.GetString(config.LogLevelFlag))
	if err != nil {
		return fmt.Errorf("could not set log level: %w", err)
	}

	logrus.SetLevel(level)
	l.log.Infof("Log level is set to %s", level)

	return nil
}

// initStorage sets the storage config to the in-memory DB or to a given Postgres DB
func (l *Launcher[T]) initStorage() (err error) {
	if viper.GetBool(config.DBInMemoryFlag) {
		l.db, err = inmemory.NewStorage()
	} else {
		l.db, err = gorm.NewStorage(gorm.WithPostgres(
			viper.GetString(config.DBHostFlag),
			viper.GetUint16(config.DBPortFlag),
			viper.GetString(config.DBUserNameFlag),
			viper.GetString(config.DBPasswordFlag),
			viper.GetString(config.DBNameFlag),
			viper.GetString(config.DBSSLModeFlag),
		))
	}
	if err != nil {
		// We could also just log the error and forward db = nil which will result in inmemory storages for each service
		// below
		return fmt.Errorf("could not create storage: %w", err)
	}

	return
}

// printClouditorHeader prints the Clouditor header for the given component
func printClouditorHeader(component string) {
	rt, _ := service.GetRuntimeInfo()

	fmt.Printf(`
           $$\                           $$\ $$\   $$\
           $$ |                          $$ |\__|  $$ |
  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 
  %s Version %s
`, component, rt.VersionString())
	fmt.Println()
}
