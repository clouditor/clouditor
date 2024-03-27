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

package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/internal/auth"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/persistence/gorm"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/server"
	"clouditor.io/clouditor/v2/server/rest"
	"clouditor.io/clouditor/v2/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	APIDefaultUserFlag               = "api-default-user"
	APIDefaultPasswordFlag           = "api-default-password"
	APIKeyPasswordFlag               = "api-key-password"
	APIKeyPathFlag                   = "api-key-path"
	APIKeySaveOnCreateFlag           = "api-key-save-on-create"
	APIgRPCPortFlag                  = "api-grpc-port"
	APIgRPCPortOrchestratorFlag      = "api-grpc-port-orchestrator"
	APIgRPCPortDiscoveryFlag         = "api-grpc-port-discovery"
	APIgRPCPortEvidenceStoreFlag     = "api-grpc-port-evidence-store"
	APIgRPCPortAssessmentFlag        = "api-grpc-port-assessment"
	APIgRPCPortEvaluationFlag        = "api-grpc-port-evaluation"
	APIHTTPPortFlag                  = "api-http-port"
	APIHTTPPortOrchestratorFlag      = "api-http-port-orchestrator"
	APIHTTPPortDiscoveryFlag         = "api-http-port-discovery"
	APIHTTPPortEvidenceStoreFlag     = "api-http-port-evidence-store"
	APIHTTPPortAssessmentFlag        = "api-http-port-assessment"
	APIHTTPPortEvaluationFlag        = "api-http-port-evaluation"
	APICORSAllowedOriginsFlags       = "api-cors-allowed-origins"
	APICORSAllowedHeadersFlags       = "api-cors-allowed-headers"
	APICORSAllowedMethodsFlags       = "api-cors-allowed-methods"
	APIJWKSURLFlag                   = "api-jwks-url"
	APIStartEmbeddedOAuth2ServerFlag = "api-start-embedded-oauth-server"
	ServiceOAuth2EndpointFlag        = "service-oauth2-token-endpoint"
	ServiceOAuth2ClientIDFlag        = "service-oauth2-client-id"
	ServiceOAuth2ClientSecretFlag    = "service-oauth2-client-secret"
	CloudServiceIDFlag               = "cloud-service-id"
	AssessmentURLFlag                = "assessment-url"
	OrchestratorURLFlag              = "orchestrator-url"
	EvidenceStoreURLFlag             = "evidence-store-url"
	DBUserNameFlag                   = "db-user-name"
	DBPasswordFlag                   = "db-password"
	DBHostFlag                       = "db-host"
	DBNameFlag                       = "db-name"
	DBPortFlag                       = "db-port"
	DBSSLModeFlag                    = "db-ssl-mode"
	DBInMemoryFlag                   = "db-in-memory"
	CreateDefaultTarget              = "target-default-create"
	DiscoveryAutoStartFlag           = "discovery-auto-start"
	DiscoveryProviderFlag            = "discovery-provider"
	DiscoveryResourceGroupFlag       = "discovery-resource-group"
	DashboardURLFlag                 = "dashboard-url"
	LogLevelFlag                     = "log-level"

	DefaultAPIDefaultUser                      = "clouditor"
	DefaultAPIDefaultPassword                  = "clouditor"
	DefaultAPIgRPCPort                  uint16 = 9090
	DefaultAPIgRPCPortOrchestrator      uint16 = 9090
	DefaultAPIgRPCPortDiscovery         uint16 = 9091
	DefaultAPIgRPCPortEvidenceStore     uint16 = 9092
	DefaultAPIgRPCPortAssessment        uint16 = 9093
	DefaultAPIgRPCPortEvaluation        uint16 = 9094
	DefaultAPIHTTPPortOrchestrator      uint16 = 8080
	DefaultAPIHTTPPortDiscovery         uint16 = 8081
	DefaultAPIHTTPPortEvidenceStore     uint16 = 8082
	DefaultAPIHTTPPortAssessment        uint16 = 8083
	DefaultAPIHTTPPortEvaluation        uint16 = 8084
	DefaultAPIStartEmbeddedOAuth2Server        = true
	DefaultServiceOAuth2Endpoint               = "http://localhost:8080/v1/auth/token"
	DefaultServiceOAuth2ClientID               = "clouditor"
	DefaultServiceOAuth2ClientSecret           = "clouditor"
	DefaultOrchestratorURL                     = "localhost:9090"
	DefaultEvidenceStoreURL                    = "localhost:9092"
	DefaultAssessmentURL                       = "localhost:9093"
	DefaultDBUserName                          = "postgres"
	DefaultDBPassword                          = "postgres"
	DefaultDBHost                              = "localhost"
	DefaultDBName                              = "postgres"
	DefaultDBPort                       uint16 = 5432
	DefaultDBSSLMode                           = "disable"
	DefaultDBInMemory                          = false
	DefaultCreateDefaultTarget                 = true
	DefaultDiscoveryAutoStart                  = false
	DefaultDiscoveryResourceGroup              = ""
	DefaultDashboardURL                        = "http://localhost:8080"
	DefaultLogLevel                            = "info"

	EnvPrefix = "CLOUDITOR"
)

func InitCobra(engineCmd *cobra.Command) *cobra.Command {
	engineCmd.Flags().String(APIDefaultUserFlag, DefaultAPIDefaultUser, "Specifies the default API username")
	engineCmd.Flags().String(APIDefaultPasswordFlag, DefaultAPIDefaultPassword, "Specifies the default API password")
	engineCmd.Flags().String(APIKeyPasswordFlag, auth.DefaultApiKeyPassword, "Specifies the password used to proctect the API private key")
	engineCmd.Flags().String(APIKeyPathFlag, auth.DefaultApiKeyPath, "Specifies the location of the API private key")
	engineCmd.Flags().Bool(APIKeySaveOnCreateFlag, auth.DefaultApiKeySaveOnCreate, "Specifies whether the API key should be saved on creation. It will only created if the default location is used.")
	engineCmd.Flags().Uint16(APIgRPCPortFlag, DefaultAPIgRPCPort, "Specifies the port used for the Clouditor gRPC API")
	engineCmd.Flags().Uint16(APIgRPCPortOrchestratorFlag, DefaultAPIgRPCPortOrchestrator, "Specifies the port used for the Orchestrator gRPC API")
	engineCmd.Flags().Uint16(APIgRPCPortDiscoveryFlag, DefaultAPIgRPCPortDiscovery, "Specifies the port used for the Discovery gRPC API")
	engineCmd.Flags().Uint16(APIgRPCPortEvidenceStoreFlag, DefaultAPIgRPCPortEvidenceStore, "Specifies the port used for the Evidence Store gRPC API")
	engineCmd.Flags().Uint16(APIgRPCPortAssessmentFlag, DefaultAPIgRPCPortAssessment, "Specifies the port used for the Assessment gRPC API")
	engineCmd.Flags().Uint16(APIgRPCPortEvaluationFlag, DefaultAPIgRPCPortEvaluation, "Specifies the port used for the Evaluation gRPC API")
	engineCmd.Flags().Uint16(APIHTTPPortFlag, rest.DefaultAPIHTTPPort, "Specifies the port used for the Clouditor HTTP API")
	engineCmd.Flags().Uint16(APIHTTPPortOrchestratorFlag, DefaultAPIHTTPPortOrchestrator, "Specifies the port used for the Orchestrator HTTP API")
	engineCmd.Flags().Uint16(APIHTTPPortDiscoveryFlag, DefaultAPIHTTPPortDiscovery, "Specifies the port used for the Discovery HTTP API")
	engineCmd.Flags().Uint16(APIHTTPPortEvidenceStoreFlag, DefaultAPIHTTPPortEvidenceStore, "Specifies the port used for the Evidence Store HTTP API")
	engineCmd.Flags().Uint16(APIHTTPPortAssessmentFlag, DefaultAPIHTTPPortAssessment, "Specifies the port used for the Assessment HTTP API")
	engineCmd.Flags().Uint16(APIHTTPPortEvaluationFlag, DefaultAPIHTTPPortEvaluation, "Specifies the port used for the Evaluation HTTP API")
	engineCmd.Flags().String(APIJWKSURLFlag, server.DefaultJWKSURL, "Specifies the JWKS URL used to verify authentication tokens in the gRPC and HTTP API")
	engineCmd.Flags().String(ServiceOAuth2EndpointFlag, DefaultServiceOAuth2Endpoint, "Specifies the OAuth 2.0 token endpoint")
	engineCmd.Flags().String(ServiceOAuth2ClientIDFlag, DefaultServiceOAuth2ClientID, "Specifies the OAuth 2.0 client ID")
	engineCmd.Flags().String(ServiceOAuth2ClientSecretFlag, DefaultServiceOAuth2ClientSecret, "Specifies the OAuth 2.0 client secret")
	engineCmd.Flags().Bool(APIStartEmbeddedOAuth2ServerFlag, DefaultAPIStartEmbeddedOAuth2Server, "Specifies whether the embedded OAuth 2.0 authorization server is started as part of the REST gateway. For production workloads, an external authorization server is recommended.")
	engineCmd.Flags().StringArray(APICORSAllowedOriginsFlags, rest.DefaultAllowedOrigins, "Specifies the origins allowed in CORS")
	engineCmd.Flags().StringArray(APICORSAllowedHeadersFlags, rest.DefaultAllowedHeaders, "Specifies the headers allowed in CORS")
	engineCmd.Flags().StringArray(APICORSAllowedMethodsFlags, rest.DefaultAllowedMethods, "Specifies the methods allowed in CORS")
	engineCmd.Flags().String(CloudServiceIDFlag, discovery.DefaultCloudServiceID, "Specifies the Cloud Service ID")
	engineCmd.Flags().String(AssessmentURLFlag, DefaultAssessmentURL, "Specifies the Assessment URL")
	engineCmd.Flags().String(OrchestratorURLFlag, DefaultOrchestratorURL, "Specifies the Orchestrator URL")
	engineCmd.Flags().String(EvidenceStoreURLFlag, DefaultEvidenceStoreURL, "Specifies the Evidence Store URL")
	engineCmd.Flags().String(DBUserNameFlag, DefaultDBUserName, "Provides user name of database")
	engineCmd.Flags().String(DBPasswordFlag, DefaultDBPassword, "Provides password of database")
	engineCmd.Flags().String(DBHostFlag, DefaultDBHost, "Provides address of database")
	engineCmd.Flags().String(DBNameFlag, DefaultDBName, "Provides name of database")
	engineCmd.Flags().Uint16(DBPortFlag, DefaultDBPort, "Provides port for database")
	engineCmd.Flags().String(DBSSLModeFlag, DefaultDBSSLMode, "The SSL mode for the database")
	engineCmd.Flags().Bool(DBInMemoryFlag, DefaultDBInMemory, "Uses an in-memory database which is not persisted at all")
	engineCmd.Flags().Bool(CreateDefaultTarget, DefaultCreateDefaultTarget, "Creates a default target cloud service if it does not exist")
	engineCmd.Flags().Bool(DiscoveryAutoStartFlag, DefaultDiscoveryAutoStart, "Automatically start the discovery when engine starts")
	engineCmd.Flags().StringSliceP(DiscoveryProviderFlag, "p", []string{}, "Providers to discover, separated by comma")
	engineCmd.Flags().String(DiscoveryResourceGroupFlag, DefaultDiscoveryResourceGroup, "Limit the scope of the discovery to a resource group (currently only used in the Azure discoverer")
	engineCmd.Flags().String(DashboardURLFlag, DefaultDashboardURL, "The URL of the Clouditor Dashboard. If the embedded server is used, a public OAuth 2.0 client based on this URL will be added")
	engineCmd.Flags().String(LogLevelFlag, DefaultLogLevel, "The default log level")

	_ = viper.BindPFlag(APIDefaultUserFlag, engineCmd.Flags().Lookup(APIDefaultUserFlag))
	_ = viper.BindPFlag(APIDefaultPasswordFlag, engineCmd.Flags().Lookup(APIDefaultPasswordFlag))
	_ = viper.BindPFlag(APIKeyPasswordFlag, engineCmd.Flags().Lookup(APIKeyPasswordFlag))
	_ = viper.BindPFlag(APIKeyPathFlag, engineCmd.Flags().Lookup(APIKeyPathFlag))
	_ = viper.BindPFlag(APIKeySaveOnCreateFlag, engineCmd.Flags().Lookup(APIKeySaveOnCreateFlag))
	_ = viper.BindPFlag(APIgRPCPortOrchestratorFlag, engineCmd.Flags().Lookup(APIgRPCPortOrchestratorFlag))
	_ = viper.BindPFlag(APIgRPCPortDiscoveryFlag, engineCmd.Flags().Lookup(APIgRPCPortDiscoveryFlag))
	_ = viper.BindPFlag(APIgRPCPortEvidenceStoreFlag, engineCmd.Flags().Lookup(APIgRPCPortEvidenceStoreFlag))
	_ = viper.BindPFlag(APIgRPCPortAssessmentFlag, engineCmd.Flags().Lookup(APIgRPCPortAssessmentFlag))
	_ = viper.BindPFlag(APIgRPCPortEvaluationFlag, engineCmd.Flags().Lookup(APIgRPCPortEvaluationFlag))
	_ = viper.BindPFlag(APIHTTPPortOrchestratorFlag, engineCmd.Flags().Lookup(APIHTTPPortOrchestratorFlag))
	_ = viper.BindPFlag(APIHTTPPortDiscoveryFlag, engineCmd.Flags().Lookup(APIHTTPPortDiscoveryFlag))
	_ = viper.BindPFlag(APIHTTPPortEvidenceStoreFlag, engineCmd.Flags().Lookup(APIHTTPPortEvidenceStoreFlag))
	_ = viper.BindPFlag(APIHTTPPortAssessmentFlag, engineCmd.Flags().Lookup(APIHTTPPortAssessmentFlag))
	_ = viper.BindPFlag(APIHTTPPortEvaluationFlag, engineCmd.Flags().Lookup(APIHTTPPortEvaluationFlag))
	_ = viper.BindPFlag(APIJWKSURLFlag, engineCmd.Flags().Lookup(APIJWKSURLFlag))
	_ = viper.BindPFlag(ServiceOAuth2EndpointFlag, engineCmd.Flags().Lookup(ServiceOAuth2EndpointFlag))
	_ = viper.BindPFlag(ServiceOAuth2ClientIDFlag, engineCmd.Flags().Lookup(ServiceOAuth2ClientIDFlag))
	_ = viper.BindPFlag(ServiceOAuth2ClientSecretFlag, engineCmd.Flags().Lookup(ServiceOAuth2ClientSecretFlag))
	_ = viper.BindPFlag(APIStartEmbeddedOAuth2ServerFlag, engineCmd.Flags().Lookup(APIStartEmbeddedOAuth2ServerFlag))
	_ = viper.BindPFlag(CloudServiceIDFlag, engineCmd.Flags().Lookup(CloudServiceIDFlag))
	_ = viper.BindPFlag(OrchestratorURLFlag, engineCmd.Flags().Lookup(OrchestratorURLFlag))
	_ = viper.BindPFlag(EvidenceStoreURLFlag, engineCmd.Flags().Lookup(EvidenceStoreURLFlag))
	_ = viper.BindPFlag(AssessmentURLFlag, engineCmd.Flags().Lookup(AssessmentURLFlag))
	_ = viper.BindPFlag(APICORSAllowedOriginsFlags, engineCmd.Flags().Lookup(APICORSAllowedOriginsFlags))
	_ = viper.BindPFlag(APICORSAllowedHeadersFlags, engineCmd.Flags().Lookup(APICORSAllowedHeadersFlags))
	_ = viper.BindPFlag(APICORSAllowedMethodsFlags, engineCmd.Flags().Lookup(APICORSAllowedMethodsFlags))
	_ = viper.BindPFlag(DBUserNameFlag, engineCmd.Flags().Lookup(DBUserNameFlag))
	_ = viper.BindPFlag(DBPasswordFlag, engineCmd.Flags().Lookup(DBPasswordFlag))
	_ = viper.BindPFlag(DBHostFlag, engineCmd.Flags().Lookup(DBHostFlag))
	_ = viper.BindPFlag(DBNameFlag, engineCmd.Flags().Lookup(DBNameFlag))
	_ = viper.BindPFlag(DBPortFlag, engineCmd.Flags().Lookup(DBPortFlag))
	_ = viper.BindPFlag(DBSSLModeFlag, engineCmd.Flags().Lookup(DBSSLModeFlag))
	_ = viper.BindPFlag(DBInMemoryFlag, engineCmd.Flags().Lookup(DBInMemoryFlag))
	_ = viper.BindPFlag(CreateDefaultTarget, engineCmd.Flags().Lookup(CreateDefaultTarget))
	_ = viper.BindPFlag(DiscoveryAutoStartFlag, engineCmd.Flags().Lookup(DiscoveryAutoStartFlag))
	_ = viper.BindPFlag(DiscoveryProviderFlag, engineCmd.Flags().Lookup(DiscoveryProviderFlag))
	_ = viper.BindPFlag(DiscoveryResourceGroupFlag, engineCmd.Flags().Lookup(DiscoveryResourceGroupFlag))
	_ = viper.BindPFlag(DashboardURLFlag, engineCmd.Flags().Lookup(DashboardURLFlag))
	_ = viper.BindPFlag(LogLevelFlag, engineCmd.Flags().Lookup(LogLevelFlag))

	return engineCmd
}

func InitConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName("clouditor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
}

// PrintClouditorHeader prints the Clouditor header for the given component
func PrintClouditorHeader(component string) {
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

// SetLogLevel sets the logrus log level
func SetLogLevel(log *logrus.Entry) (*logrus.Entry, error) {
	level, err := logrus.ParseLevel(viper.GetString(LogLevelFlag))
	if err != nil {
		return nil, fmt.Errorf("could not set log level: %w", err)
	}

	logrus.SetLevel(level)
	log.Infof("Log level is set to %s", level)

	return log, nil
}

// SetStorage sets the storage config to the in-memory DB or to a given Postgres DB
func SetStorage() (db persistence.Storage, err error) {
	if viper.GetBool(DBInMemoryFlag) {
		db, err = inmemory.NewStorage()
	} else {
		db, err = gorm.NewStorage(gorm.WithPostgres(
			viper.GetString(DBHostFlag),
			viper.GetUint16(DBPortFlag),
			viper.GetString(DBUserNameFlag),
			viper.GetString(DBPasswordFlag),
			viper.GetString(DBNameFlag),
			viper.GetString(DBSSLModeFlag),
		))
	}
	if err != nil {
		// We could also just log the error and forward db = nil which will result in inmemory storages for each service
		// below
		return nil, fmt.Errorf("could not create storage: %w", err)
	}

	return
}

// StartServer starts the gRPC server and the corresponding gRPC-HTTP gateway with the given gRPC Server Options
func StartServer(log *logrus.Entry, grpcOpts ...server.StartGRPCServerOption) (srv *server.Server, err error) {
	var (
		grpcPort uint16
		httpPort uint16
		restOpts []rest.ServerConfigOption
	)

	grpcPort = viper.GetUint16(APIgRPCPortOrchestratorFlag)
	httpPort = viper.GetUint16(APIHTTPPortOrchestratorFlag)

	restOpts = []rest.ServerConfigOption{
		rest.WithAllowedOrigins(viper.GetStringSlice(APICORSAllowedOriginsFlags)),
		rest.WithAllowedHeaders(viper.GetStringSlice(APICORSAllowedHeadersFlags)),
		rest.WithAllowedMethods(viper.GetStringSlice(APICORSAllowedMethodsFlags)),
	}

	log.Infof("Starting gRPC endpoint on :%d", grpcPort)

	// Start the gRPC server
	_, srv, err = server.StartGRPCServer(
		fmt.Sprintf("0.0.0.0:%d", grpcPort),
		grpcOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to serve gRPC endpoint: %w", err)
	}

	// Start the gRPC-HTTP gateway
	err = rest.RunServer(context.Background(),
		grpcPort,
		httpPort,
		restOpts...,
	)
	if err != nil && err != http.ErrServerClosed {
		return nil, fmt.Errorf("failed to serve gRPC-HTTP gateway: %v", err)
	}

	log.Infof("Stopping gRPC endpoint")
	srv.Stop()

	return
}
