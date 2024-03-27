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

package commands

import (
	"clouditor.io/clouditor/v2/internal/auth"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(config.APIDefaultUserFlag, config.DefaultAPIDefaultUser, "Specifies the default API username")
	cmd.PersistentFlags().String(config.APIDefaultPasswordFlag, config.DefaultAPIDefaultPassword, "Specifies the default API password")
	cmd.PersistentFlags().String(config.APIKeyPasswordFlag, auth.DefaultApiKeyPassword, "Specifies the password used to protect the API private key")
	cmd.PersistentFlags().String(config.APIKeyPathFlag, auth.DefaultApiKeyPath, "Specifies the location of the API private key")
	cmd.PersistentFlags().Bool(config.APIKeySaveOnCreateFlag, auth.DefaultApiKeySaveOnCreate, "Specifies whether the API key should be saved on creation. It will only created if the default location is used.")
	cmd.PersistentFlags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPort, "Specifies the port used for the Clouditor gRPC API")
	cmd.PersistentFlags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPort, "Specifies the port used for the Clouditor HTTP API")
	cmd.PersistentFlags().String(config.APIJWKSURLFlag, server.DefaultJWKSURL, "Specifies the JWKS URL used to verify authentication tokens in the gRPC and HTTP API")
	cmd.PersistentFlags().StringArray(config.APICORSAllowedOriginsFlags, config.DefaultAllowedOrigins, "Specifies the origins allowed in CORS")
	cmd.PersistentFlags().StringArray(config.APICORSAllowedHeadersFlags, config.DefaultAllowedHeaders, "Specifies the headers allowed in CORS")
	cmd.PersistentFlags().StringArray(config.APICORSAllowedMethodsFlags, config.DefaultAllowedMethods, "Specifies the methods allowed in CORS")
	cmd.PersistentFlags().String(config.DashboardURLFlag, config.DefaultDashboardURL, "The URL of the Clouditor Dashboard. If the embedded server is used, a public OAuth 2.0 client based on this URL will be added")
	cmd.PersistentFlags().String(config.DBUserNameFlag, config.DefaultDBUserName, "Provides user name of database")
	cmd.PersistentFlags().String(config.DBPasswordFlag, config.DefaultDBPassword, "Provides password of database")
	cmd.PersistentFlags().String(config.DBHostFlag, config.DefaultDBHost, "Provides address of database")
	cmd.PersistentFlags().String(config.DBNameFlag, config.DefaultDBName, "Provides name of database")
	cmd.PersistentFlags().Uint16(config.DBPortFlag, config.DefaultDBPort, "Provides port for database")
	cmd.PersistentFlags().String(config.DBSSLModeFlag, config.DefaultDBSSLMode, "The SSL mode for the database")
	cmd.PersistentFlags().Bool(config.DBInMemoryFlag, config.DefaultDBInMemory, "Uses an in-memory database which is not persisted at all")
	cmd.PersistentFlags().String(config.ServiceOAuth2EndpointFlag, config.DefaultServiceOAuth2Endpoint, "Specifies the OAuth 2.0 token endpoint")
	cmd.PersistentFlags().String(config.ServiceOAuth2ClientIDFlag, config.DefaultServiceOAuth2ClientID, "Specifies the OAuth 2.0 client ID")
	cmd.PersistentFlags().String(config.ServiceOAuth2ClientSecretFlag, config.DefaultServiceOAuth2ClientSecret, "Specifies the OAuth 2.0 client secret")
	cmd.PersistentFlags().String(config.LogLevelFlag, config.DefaultLogLevel, "The default log level")

	_ = viper.BindPFlag(config.APIDefaultUserFlag, cmd.PersistentFlags().Lookup(config.APIDefaultUserFlag))
	_ = viper.BindPFlag(config.APIDefaultPasswordFlag, cmd.PersistentFlags().Lookup(config.APIDefaultPasswordFlag))
	_ = viper.BindPFlag(config.APIKeyPasswordFlag, cmd.PersistentFlags().Lookup(config.APIKeyPasswordFlag))
	_ = viper.BindPFlag(config.APIKeyPathFlag, cmd.PersistentFlags().Lookup(config.APIKeyPathFlag))
	_ = viper.BindPFlag(config.APIKeySaveOnCreateFlag, cmd.PersistentFlags().Lookup(config.APIKeySaveOnCreateFlag))
	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.PersistentFlags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.PersistentFlags().Lookup(config.APIHTTPPortFlag))
	_ = viper.BindPFlag(config.APIJWKSURLFlag, cmd.PersistentFlags().Lookup(config.APIJWKSURLFlag))
	_ = viper.BindPFlag(config.APICORSAllowedOriginsFlags, cmd.PersistentFlags().Lookup(config.APICORSAllowedOriginsFlags))
	_ = viper.BindPFlag(config.APICORSAllowedHeadersFlags, cmd.PersistentFlags().Lookup(config.APICORSAllowedHeadersFlags))
	_ = viper.BindPFlag(config.APICORSAllowedMethodsFlags, cmd.PersistentFlags().Lookup(config.APICORSAllowedMethodsFlags))
	_ = viper.BindPFlag(config.DashboardURLFlag, cmd.PersistentFlags().Lookup(config.DashboardURLFlag))
	_ = viper.BindPFlag(config.DBUserNameFlag, cmd.PersistentFlags().Lookup(config.DBUserNameFlag))
	_ = viper.BindPFlag(config.DBPasswordFlag, cmd.PersistentFlags().Lookup(config.DBPasswordFlag))
	_ = viper.BindPFlag(config.DBHostFlag, cmd.PersistentFlags().Lookup(config.DBHostFlag))
	_ = viper.BindPFlag(config.DBNameFlag, cmd.PersistentFlags().Lookup(config.DBNameFlag))
	_ = viper.BindPFlag(config.DBPortFlag, cmd.PersistentFlags().Lookup(config.DBPortFlag))
	_ = viper.BindPFlag(config.DBSSLModeFlag, cmd.PersistentFlags().Lookup(config.DBSSLModeFlag))
	_ = viper.BindPFlag(config.DBInMemoryFlag, cmd.PersistentFlags().Lookup(config.DBInMemoryFlag))
	_ = viper.BindPFlag(config.ServiceOAuth2EndpointFlag, cmd.PersistentFlags().Lookup(config.ServiceOAuth2EndpointFlag))
	_ = viper.BindPFlag(config.ServiceOAuth2ClientIDFlag, cmd.PersistentFlags().Lookup(config.ServiceOAuth2ClientIDFlag))
	_ = viper.BindPFlag(config.ServiceOAuth2ClientSecretFlag, cmd.PersistentFlags().Lookup(config.ServiceOAuth2ClientSecretFlag))
	_ = viper.BindPFlag(config.LogLevelFlag, cmd.PersistentFlags().Lookup(config.LogLevelFlag))
}
