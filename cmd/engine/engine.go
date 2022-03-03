// Copyright 2016-2020 Fraunhofer AISEC
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

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	cli_discovery "clouditor.io/clouditor/cli/commands/service/discovery"
	"clouditor.io/clouditor/logging/formatter"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/persistence/gorm"
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/rest"
	"clouditor.io/clouditor/service"
	service_assessment "clouditor.io/clouditor/service/assessment"
	service_auth "clouditor.io/clouditor/service/auth"
	service_discovery "clouditor.io/clouditor/service/discovery"
	service_evidenceStore "clouditor.io/clouditor/service/evidence"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	APIDefaultUserFlag         = "api-default-user"
	APIDefaultPasswordFlag     = "api-default-password"
	APIKeyPasswordFlag         = "api-key-password"
	APIKeyPathFlag             = "api-key-path"
	APIKeySaveOnCreateFlag     = "api-key-save-on-create"
	APIgRPCPortFlag            = "api-grpc-port"
	APIHTTPPortFlag            = "api-http-port"
	APICORSAllowedOriginsFlags = "api-cors-allowed-origins"
	APICORSAllowedHeadersFlags = "api-cors-allowed-headers"
	APICORSAllowedMethodsFlags = "api-cors-allowed-methods"
	APIJWKSURLFlag             = "api-jwks-url"
	DBUserNameFlag             = "db-user-name"
	DBPasswordFlag             = "db-password"
	DBHostFlag                 = "db-host"
	DBNameFlag                 = "db-name"
	DBPortFlag                 = "db-port"
	DBInMemoryFlag             = "db-in-memory"
	CreateDefaultTarget        = "target-default-create"

	DefaultAPIDefaultUser      = "clouditor"
	DefaultAPIDefaultPassword  = "clouditor"
	DefaultAPIgRPCPort         = 9090
	DefaultDBUserName          = "postgres"
	DefaultDBPassword          = "postgres"
	DefaultDBHost              = "localhost"
	DefaultDBName              = "postgres"
	DefaultDBPort              = 5432
	DefaultDBInMemory          = false
	DefaultCreateDefaultTarget = true

	EnvPrefix = "CLOUDITOR"
)

var (
	server               *grpc.Server
	authService          *service_auth.Service
	discoveryService     discovery.DiscoveryServer
	orchestratorService  *service_orchestrator.Service
	assessmentService    assessment.AssessmentServer
	evidenceStoreService evidence.EvidenceStoreServer
	db                   persistence.Storage

	log *logrus.Entry
)

var engineCmd = &cobra.Command{
	Use:   "engine",
	Short: "engine launches the Clouditor Engine",
	Long:  "Clouditor Engine is the main component of Clouditor. It is an all-in-one solution of several microservices, which also can be started individually.",
	RunE:  doCmd,
}

func init() {
	log = logrus.WithField("component", "grpc")
	log.Logger.Formatter = formatter.CapitalizeFormatter{Formatter: &logrus.TextFormatter{ForceColors: true}}

	cobra.OnInitialize(initConfig)

	engineCmd.Flags().String(APIDefaultUserFlag, DefaultAPIDefaultUser, "Specifies the default API username")
	engineCmd.Flags().String(APIDefaultPasswordFlag, DefaultAPIDefaultPassword, "Specifies the default API password")
	engineCmd.Flags().String(APIKeyPasswordFlag, service_auth.DefaultApiKeyPassword, "Specifies the password used to proctect the API private key")
	engineCmd.Flags().String(APIKeyPathFlag, service_auth.DefaultApiKeyPath, "Specifies the location of the API private key")
	engineCmd.Flags().Bool(APIKeySaveOnCreateFlag, service_auth.DefaultApiKeySaveOnCreate, "Specifies whether the API key should be saved on creation. It will only created if the default location is used.")
	engineCmd.Flags().Int16(APIgRPCPortFlag, DefaultAPIgRPCPort, "Specifies the port used for the gRPC API")
	engineCmd.Flags().Int16(APIHTTPPortFlag, rest.DefaultAPIHTTPPort, "Specifies the port used for the HTTP API")
	engineCmd.Flags().String(APIJWKSURLFlag, service.DefaultJWKSURL, "Specifies the JWKS URL used to verify authentication tokens in the gRPC and HTTP API")
	engineCmd.Flags().StringArray(APICORSAllowedOriginsFlags, rest.DefaultAllowedOrigins, "Specifies the origins allowed in CORS")
	engineCmd.Flags().StringArray(APICORSAllowedHeadersFlags, rest.DefaultAllowedHeaders, "Specifies the headers allowed in CORS")
	engineCmd.Flags().StringArray(APICORSAllowedMethodsFlags, rest.DefaultAllowedMethods, "Specifies the methods allowed in CORS")
	engineCmd.Flags().String(DBUserNameFlag, DefaultDBUserName, "Provides user name of database")
	engineCmd.Flags().String(DBPasswordFlag, DefaultDBPassword, "Provides password of database")
	engineCmd.Flags().String(DBHostFlag, DefaultDBHost, "Provides address of database")
	engineCmd.Flags().String(DBNameFlag, DefaultDBName, "Provides name of database")
	engineCmd.Flags().Int16(DBPortFlag, DefaultDBPort, "Provides port for database")
	engineCmd.Flags().Bool(DBInMemoryFlag, DefaultDBInMemory, "Uses an in-memory database which is not persisted at all")
	engineCmd.Flags().Bool(CreateDefaultTarget, DefaultCreateDefaultTarget, "Creates a default target cloud service if it does not exist")

	_ = viper.BindPFlag(APIDefaultUserFlag, engineCmd.Flags().Lookup(APIDefaultUserFlag))
	_ = viper.BindPFlag(APIDefaultPasswordFlag, engineCmd.Flags().Lookup(APIDefaultPasswordFlag))
	_ = viper.BindPFlag(APIKeyPasswordFlag, engineCmd.Flags().Lookup(APIKeyPasswordFlag))
	_ = viper.BindPFlag(APIKeyPathFlag, engineCmd.Flags().Lookup(APIKeyPathFlag))
	_ = viper.BindPFlag(APIKeySaveOnCreateFlag, engineCmd.Flags().Lookup(APIKeySaveOnCreateFlag))
	_ = viper.BindPFlag(APIgRPCPortFlag, engineCmd.Flags().Lookup(APIgRPCPortFlag))
	_ = viper.BindPFlag(APIHTTPPortFlag, engineCmd.Flags().Lookup(APIHTTPPortFlag))
	_ = viper.BindPFlag(APIJWKSURLFlag, engineCmd.Flags().Lookup(APIJWKSURLFlag))
	_ = viper.BindPFlag(APICORSAllowedOriginsFlags, engineCmd.Flags().Lookup(APICORSAllowedOriginsFlags))
	_ = viper.BindPFlag(APICORSAllowedHeadersFlags, engineCmd.Flags().Lookup(APICORSAllowedHeadersFlags))
	_ = viper.BindPFlag(APICORSAllowedMethodsFlags, engineCmd.Flags().Lookup(APICORSAllowedMethodsFlags))
	_ = viper.BindPFlag(DBUserNameFlag, engineCmd.Flags().Lookup(DBUserNameFlag))
	_ = viper.BindPFlag(DBPasswordFlag, engineCmd.Flags().Lookup(DBPasswordFlag))
	_ = viper.BindPFlag(DBHostFlag, engineCmd.Flags().Lookup(DBHostFlag))
	_ = viper.BindPFlag(DBNameFlag, engineCmd.Flags().Lookup(DBNameFlag))
	_ = viper.BindPFlag(DBPortFlag, engineCmd.Flags().Lookup(DBPortFlag))
	_ = viper.BindPFlag(DBInMemoryFlag, engineCmd.Flags().Lookup(DBInMemoryFlag))
	_ = viper.BindPFlag(CreateDefaultTarget, engineCmd.Flags().Lookup(CreateDefaultTarget))
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName("clouditor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
}

func doCmd(_ *cobra.Command, _ []string) (err error) {
	log.Info("Welcome to new Clouditor 2.0")

	fmt.Println(`
           $$\                           $$\ $$\   $$\
           $$ |                          $$ |\__|  $$ |
  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 `)

	if viper.GetBool(DBInMemoryFlag) {
		db, err = inmemory.NewStorage()
	} else {
		db, err = gorm.NewStorage(gorm.WithPostgres(viper.GetString(DBHostFlag), int16(viper.GetInt(DBPortFlag))))
	}
	if err != nil {
		// We could also just log the error and forward db = nil which will result in inmemory storages
		// for each service below
		return fmt.Errorf("could not create storage: %w", err)
	}

	authService = service_auth.NewService(
		service_auth.WithStorage(db),
		service_auth.WithApiKeyPassword(viper.GetString(APIKeyPasswordFlag)),
		service_auth.WithApiKeyPath(viper.GetString(APIKeyPathFlag)),
		service_auth.WithApiKeySaveOnCreate(viper.GetBool(APIKeySaveOnCreateFlag)),
	)
	discoveryService = service_discovery.NewService(
		service_discovery.WithInternalAuthorizer(
			api.DefaultInternalAuthorizerAddress,
			viper.GetString(APIDefaultUserFlag),
			viper.GetString(APIDefaultPasswordFlag),
		),
		service_discovery.WithDiscoverers([]string{viper.GetString(cli_discovery.DISCOVERER)}),
	)
	orchestratorService = service_orchestrator.NewService(service_orchestrator.WithStorage(db))
	assessmentService = service_assessment.NewService(
		service_assessment.WithInternalAuthorizer(
			api.DefaultInternalAuthorizerAddress,
			viper.GetString(APIDefaultUserFlag),
			viper.GetString(APIDefaultPasswordFlag),
		),
	)
	evidenceStoreService = service_evidenceStore.NewService()

	// It is possible to register hook functions for the orchestrator, evidenceStore and assessment service.
	// The hook functions in orchestrator are implemented in StoreAssessmentResult(s)
	// The hook functions in evidenceStore are implemented in StoreEvidence(s)
	// The hook functions in assessment are implemented in AssessEvidence(s)

	// orchestratorService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {})
	// evidenceStoreService.RegisterEvidenceHook(func(result *evidence.Evidence, err error) {})
	// assessmentService.RegisterAssessmentResultHook(func(result *assessment.AssessmentResult, err error) {}

	err = authService.CreateDefaultUser(viper.GetString(APIDefaultUserFlag), viper.GetString(APIDefaultPasswordFlag))
	if err != nil {
		log.Errorf("Could not create default user: %v", err)
	}

	if viper.GetBool(CreateDefaultTarget) {
		_, err := orchestratorService.CreateDefaultTargetCloudService()
		if err != nil {
			log.Errorf("could not register default target cloud service: %v", err)
		}
	}

	grpcPort := viper.GetInt(APIgRPCPortFlag)
	httpPort := viper.GetInt(APIHTTPPortFlag)

	grpcLogger := logrus.New()
	grpcLogger.Formatter = &formatter.GRPCFormatter{TextFormatter: logrus.TextFormatter{ForceColors: true}}
	grpcLoggerEntry := grpcLogger.WithField("component", "grpc")

	// disabling the grpc log itself, because it will log everything on INFO, whereas DEBUG would be more
	// appropriate
	// grpc_logrus.ReplaceGrpcLogger(grpcLoggerEntry)

	// create a new socket for gRPC communication
	sock, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Errorf("could not listen: %v", err)
	}

	authConfig := service.ConfigureAuth(service.WithJWKSURL(viper.GetString(APIJWKSURLFlag)))
	defer authConfig.Jwks.EndBackground()

	server = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(grpcLoggerEntry),
			grpc_auth.UnaryServerInterceptor(authConfig.AuthFunc),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(grpcLoggerEntry),
			grpc_auth.StreamServerInterceptor(authConfig.AuthFunc),
		))
	auth.RegisterAuthenticationServer(server, authService)
	discovery.RegisterDiscoveryServer(server, discoveryService)
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)
	assessment.RegisterAssessmentServer(server, assessmentService)
	evidence.RegisterEvidenceStoreServer(server, evidenceStoreService)

	// enable reflection, primary for testing in early stages
	reflection.Register(server)

	// start the gRPC-HTTP gateway
	go func() {
		err = rest.RunServer(context.Background(),
			grpcPort,
			httpPort,
			rest.WithAllowedOrigins(viper.GetStringSlice(APICORSAllowedOriginsFlags)),
			rest.WithAllowedHeaders(viper.GetStringSlice(APICORSAllowedHeadersFlags)),
			rest.WithAllowedMethods(viper.GetStringSlice(APICORSAllowedMethodsFlags)),
		)
		if errors.Is(err, http.ErrServerClosed) {
			// ToDo(oxisto): deepsource anti-pattern: calls to os.Exit only in main() or init() functions
			os.Exit(0)
			return
		}

		if err != nil {
			// ToDo(oxisto): deepsource anti-pattern: calls to log.Fatalf only in main() or init() functions
			log.Fatalf("failed to serve gRPC-HTTP gateway: %v", err)
		}
	}()

	log.Infof("Starting gRPC endpoint on :%d", grpcPort)

	// serve the gRPC socket
	if err := server.Serve(sock); err != nil {
		log.Infof("Failed to serve gRPC endpoint: %s", err)
		return err
	}

	return nil
}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
