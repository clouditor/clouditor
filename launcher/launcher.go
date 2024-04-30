// Copyright 2024 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package launcher

import (
	"context"
	"fmt"
	"net/http"

	commands_login "clouditor.io/clouditor/v2/cli/commands/login"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/logging/formatter"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/login"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Launcher struct {
	name     string
	srv      *server.Server
	db       persistence.Storage
	log      *logrus.Entry
	grpcOpts []server.StartGRPCServerOption
	services []service.Service
}

type NewServiceFunc[T service.Service] func(opts ...service.Option[T]) T
type WithStorageFunc[T service.Service] func(db persistence.Storage) service.Option[T]
type ServiceInitFunc[T service.Service] func(svc T) ([]server.StartGRPCServerOption, error)

func NewLauncher(name string, specs ...ServiceSpec) (l *Launcher, err error) {
	l = &Launcher{
		name: name,
	}

	// Print Clouditor header with the used Clouditor version
	printClouditorHeader(fmt.Sprintf("Clouditor %s Service", cases.Title(language.English).String(l.name)))

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

	// Create the services out of the service specs
	for _, spec := range specs {
		// Create the service and gather the gRPC server options

		svc, grpcOpts, err := spec.NewService(l.db)
		if err != nil {
			return nil, err
		}

		// Append the gRPC server options
		l.grpcOpts = append(l.grpcOpts, grpcOpts...)

		// Add the service to the list of our managed services
		l.services = append(l.services, svc)
	}

	return
}

// Launch starts the gRPC server and the corresponding gRPC-HTTP gateway with the given gRPC server Options
func (l *Launcher) Launch() (err error) {
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
		server.WithServices(l.services),
	}

	// Append launch-specific ones
	grpcOpts = append(grpcOpts, l.grpcOpts...)

	grpcPort = viper.GetUint16(config.APIgRPCPortFlag)
	httpPort = viper.GetUint16(config.APIHTTPPortFlag)

	restOpts = []rest.ServerConfigOption{
		rest.WithAllowedOrigins(viper.GetStringSlice(config.APICORSAllowedOriginsFlags)),
		rest.WithAllowedHeaders(viper.GetStringSlice(config.APICORSAllowedHeadersFlags)),
		rest.WithAllowedMethods(viper.GetStringSlice(config.APICORSAllowedMethodsFlags)),
	}

	// Let's check, if we are using our embedded OAuth 2.0 server, which we need to start (using additional arguments to
	// our existing REST gateway). In a production scenario the usage of a dedicated (external) OAuth 2.0 server is
	// recommended. In order to configure the external server, the flags ServiceOAuth2EndpointFlag and APIJWKSURLFlag
	// can be used.
	if viper.GetBool(config.APIStartEmbeddedOAuth2ServerFlag) {
		restOpts = append(restOpts,
			rest.WithEmbeddedOAuth2Server(
				viper.GetString(config.APIKeyPathFlag),
				viper.GetString(config.APIKeyPasswordFlag),
				viper.GetBool(config.APIKeySaveOnCreateFlag),
				// Create a public client for our CLI
				oauth2.WithClient(
					commands_login.DefaultClientID,
					"",
					commands_login.DefaultCallback,
				),
				// Create a public client for our dashboard
				oauth2.WithClient(
					"dashboard",
					"",
					fmt.Sprintf("%s/callback", viper.GetString(config.DashboardURLFlag)),
				),
				// Create a confidential client with default credentials for our services
				oauth2.WithClient(
					viper.GetString(config.ServiceOAuth2ClientIDFlag),
					viper.GetString(config.ServiceOAuth2ClientIDFlag),
					"",
				),
				// Createa a default user for logging in
				login.WithLoginPage(
					login.WithUser(
						viper.GetString(config.APIDefaultUserFlag),
						viper.GetString(config.APIDefaultPasswordFlag),
					),
					login.WithBaseURL("/v1/auth"),
				),
			),
		)
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

	// Do any post-start initialization of the services
	for _, svc := range l.services {
		svc.Init()
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

	for _, svc := range l.services {
		l.log.Infof("Stopping %T service", svc)
		svc.Shutdown()
	}

	l.log.Infof("Stopping gRPC endpoint")
	l.srv.Stop()

	return nil
}

// initLogging initializes the logging
func (l *Launcher) initLogging() error {
	l.log = logrus.WithField("launcher", l.name)
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
func (l *Launcher) initStorage() (err error) {
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
