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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	APIDefaultUserFlag                   = "api-default-user"
	APIDefaultPasswordFlag               = "api-default-password"
	APIKeyPasswordFlag                   = "api-key-password"
	APIKeyPathFlag                       = "api-key-path"
	APIKeySaveOnCreateFlag               = "api-key-save-on-create"
	APIgRPCPortFlag                      = "api-grpc-port"
	APIHTTPPortFlag                      = "api-http-port"
	APICORSAllowedOriginsFlags           = "api-cors-allowed-origins"
	APICORSAllowedHeadersFlags           = "api-cors-allowed-headers"
	APICORSAllowedMethodsFlags           = "api-cors-allowed-methods"
	APIJWKSURLFlag                       = "api-jwks-url"
	APIStartEmbeddedOAuth2ServerFlag     = "api-start-embedded-oauth-server"
	ServiceOAuth2EndpointFlag            = "service-oauth2-token-endpoint"
	ServiceOAuth2ClientIDFlag            = "service-oauth2-client-id"
	ServiceOAuth2ClientSecretFlag        = "service-oauth2-client-secret"
	CertificationTargetIDFlag            = "cloud-service-id"
	AssessmentURLFlag                    = "assessment-url"
	OrchestratorURLFlag                  = "orchestrator-url"
	EvidenceStoreURLFlag                 = "evidence-store-url"
	DBUserNameFlag                       = "db-user-name"
	DBPasswordFlag                       = "db-password"
	DBHostFlag                           = "db-host"
	DBNameFlag                           = "db-name"
	DBPortFlag                           = "db-port"
	DBSSLModeFlag                        = "db-ssl-mode"
	DBInMemoryFlag                       = "db-in-memory"
	CreateDefaultCertificationTargetFlag = "create-default-certification-target"
	DiscoveryAutoStartFlag               = "discovery-auto-start"
	DiscoveryProviderFlag                = "discovery-provider"
	DiscoveryResourceGroupFlag           = "discovery-resource-group"
	DiscoveryCSAFDomainFlag              = "discovery-csaf-domain"
	DashboardCallbackURLFlag             = "dashboard-callback-url"
	LogLevelFlag                         = "log-level"

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
	DefaultEvidenceStoreURL                    = "localhost:9090"
	DefaultAssessmentURL                       = "localhost:9090"
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
	DefaultCSAFDomain                          = ""
	DefaultDashboardCallbackURL                = "http://localhost:8080/callback"
	DefaultLogLevel                            = "info"

	EnvPrefix = "CLOUDITOR"
)

var (
	// DefaultAllowedOrigins contains a nil slice, as per default, no origins are allowed.
	DefaultAllowedOrigins []string = nil

	// DefaultAllowedHeaders contains sensible defaults for the Access-Control-Allow-Headers header.
	// Please adjust accordingly in production using WithAllowedHeaders.
	DefaultAllowedHeaders = []string{"Content-Type", "Accept", "Authorization"}

	// DefaultAllowedMethods contains sensible defaults for the Access-Control-Allow-Methods header.
	// Please adjust accordingly in production using WithAllowedMethods.
	DefaultAllowedMethods = []string{"GET", "POST", "PUT", "DELETE"}

	// DefaultAPIHTTPPort specifies the default port for the REST API.
	DefaultAPIHTTPPort uint16 = 8080
)

const (
	// DefaultCertificationTargetID is the default service ID. Currently, our discoverers have no way to differentiate between different
	// services, but we need this feature in the future. This serves as a default to already prepare the necessary
	// structures for this feature.
	DefaultCertificationTargetID = "00000000-0000-0000-0000-000000000000"

	// DefaultEvidenceCollectorToolID is the default evidence collector tool ID.
	DefaultEvidenceCollectorToolID = "Clouditor Evidences Collection"
)

func init() {
	cobra.OnInitialize(InitConfig)
}

func InitCobra(engineCmd *cobra.Command) {
	engineCmd.Flags().Bool(APIStartEmbeddedOAuth2ServerFlag, DefaultAPIStartEmbeddedOAuth2Server, "Specifies whether the embedded OAuth 2.0 authorization server is started as part of the REST gateway. For production workloads, an external authorization server is recommended.")

	_ = viper.BindPFlag(APIStartEmbeddedOAuth2ServerFlag, engineCmd.Flags().Lookup(APIStartEmbeddedOAuth2ServerFlag))
}

func InitConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.SetConfigName("clouditor")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}

// ClientCredentials configures the OAuth 2.0 client credentials for a service
func ClientCredentials() *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     viper.GetString(ServiceOAuth2ClientIDFlag),
		ClientSecret: viper.GetString(ServiceOAuth2ClientSecretFlag),
		TokenURL:     viper.GetString(ServiceOAuth2EndpointFlag),
	}
}
