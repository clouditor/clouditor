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
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	APIDefaultUserFlag     = "api-default-user"
	APIDefaultPasswordFlag = "api-default-password"
	APISecretFlag          = "api-secret"
	APIgRPCPortFlag        = "api-grpc-port"
	APIHTTPPortFlag        = "api-http-port"
	DBUserNameFlag         = "db-user-name"
	DBPasswordFlag         = "db-password"
	DBHostFlag             = "db-host"
	DBNameFlag             = "db-name"
	DBPortFlag             = "db-port"
	DBInMemoryFlag         = "db-in-memory"

	DefaultAPIDefaultUser     = "clouditor"
	DefaultAPIDefaultPassword = "clouditor"
	DefaultAPISecret          = "changeme"
	DefaultAPIgRPCPort        = 9090
	DefaultAPIHTTPPort        = 8080
	DefaultDBUserName         = "postgres"
	DefaultDBPassword         = "postgres"
	DefaultDBHost             = "localhost"
	DefaultDBName             = "postgres"
	DefaultDBPort             = 5432
	DefaultDBInMemory         = false

	EnvPrefix = "CLOUDITOR"
)

var log *logrus.Entry

var collectorCmd = &cobra.Command{
	Use:   "collector",
	Short: "collector launches the Clouditor Collector",
	Long:  "Clouditor Collector is in charge of collecting evidence from various cloud providers",
	RunE:  doCmd,
}

func init() {
	log = logrus.WithField("component", "grpc")

	cobra.OnInitialize(initConfig)

	// TODO(oxisto): We share these flags with the engine command?
	collectorCmd.Flags().String(APIDefaultUserFlag, DefaultAPIDefaultUser, "Specifies the default API username")
	collectorCmd.Flags().String(APIDefaultPasswordFlag, DefaultAPIDefaultPassword, "Specifies the default API password")
	collectorCmd.Flags().String(APISecretFlag, DefaultAPISecret, "Specifies the secret used by API tokens")
	collectorCmd.Flags().Int16(APIgRPCPortFlag, DefaultAPIgRPCPort, "Specifies the port used for the gRPC API")
	collectorCmd.Flags().Int16(APIHTTPPortFlag, DefaultAPIHTTPPort, "Specifies the port used for the HTTP API")

	_ = viper.BindPFlag(APIDefaultUserFlag, collectorCmd.Flags().Lookup(APIDefaultUserFlag))
	_ = viper.BindPFlag(APIDefaultPasswordFlag, collectorCmd.Flags().Lookup(APIDefaultPasswordFlag))
	_ = viper.BindPFlag(APISecretFlag, collectorCmd.Flags().Lookup(APISecretFlag))
	_ = viper.BindPFlag(APIgRPCPortFlag, collectorCmd.Flags().Lookup(APIgRPCPortFlag))
	_ = viper.BindPFlag(APIHTTPPortFlag, collectorCmd.Flags().Lookup(APIHTTPPortFlag))
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
	return nil
}

func main() {
	if err := collectorCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
