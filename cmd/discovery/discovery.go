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

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/launcher"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var engineCmd = &cobra.Command{
	Use:   "discovery",
	Short: "discovery launches the Clouditor Discovery Service",
	Long:  "Discovery is a component of the Clouditor and starts the Discovery Service.",
	RunE:  doCmd,
}

func init() {
	config.InitCobra(engineCmd)
}

func doCmd(cmd *cobra.Command, _ []string) (err error) {
	var providers []string

	// If no CSPs for discovering are given, take all implemented discoverers
	if len(viper.GetStringSlice(config.DiscoveryProviderFlag)) == 0 {
		providers = []string{service_discovery.ProviderAWS, service_discovery.ProviderAzure, service_discovery.ProviderK8S}
	} else {
		providers = viper.GetStringSlice(config.DiscoveryProviderFlag)
	}

	l, err := launcher.NewLauncher(
		cmd.Use,
		service_discovery.DefaultServiceSpec,
	)
	if err != nil {
		return err
	}

	// Start the gRPC server and the corresponding gRPC-HTTP gateway
	return l.Launch()

}

func main() {
	if err := engineCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
