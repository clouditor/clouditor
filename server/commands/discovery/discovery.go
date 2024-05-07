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

package discovery

import (
	"fmt"
	"strconv"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/service/discovery"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDiscoveryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discovery",
		Short: "Starts a server which contains the Clouditor Discovery Service",
		Long:  "This command starts a Clouditor Discovery service",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				discovery.DefaultServiceSpec(),
			)
			if err != nil {
				return err
			}

			return l.Launch()
		},
	}

	BindFlags(cmd)

	return cmd
}

func BindFlags(cmd *cobra.Command) {
	// Set the AssessmentURLFlag default value to the default assessment gRPC port, e.g., "localhost:9093"
	if cmd.Flag(config.AssessmentURLFlag) == nil {
		cmd.Flags().String(config.AssessmentURLFlag, fmt.Sprintf("localhost:%s", strconv.FormatUint(uint64(config.DefaultAPIgRPCPortAssessment), 10)), "Specifies the Assessment URL")
	}
	cmd.Flags().String(config.CloudServiceIDFlag, config.DefaultCloudServiceID, "Specifies the Cloud Service ID")
	cmd.Flags().Bool(config.DiscoveryAutoStartFlag, config.DefaultDiscoveryAutoStart, "Automatically start the discovery when engine starts")
	cmd.Flags().StringSliceP(config.DiscoveryProviderFlag, "p", []string{}, "Providers to discover, separated by comma")
	cmd.Flags().String(config.DiscoveryResourceGroupFlag, config.DefaultDiscoveryResourceGroup, "Limit the scope of the discovery to a resource group (currently only used in the Azure discoverer")
	cmd.Flags().String(config.DiscoveryCSAFDomainFlag, config.DefaultCSAFDomain, "The domain to look for a CSAF provider, if the CSAF discovery is enabled")
	if cmd.Flag(config.APIgRPCPortFlag) == nil {
		cmd.Flags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPortDiscovery, "Specifies the port used for the Clouditor gRPC API")
	}
	if cmd.Flag(config.APIHTTPPortFlag) == nil {
		cmd.Flags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPortDiscovery, "Specifies the port used for the Clouditor HTTP API")
	}

	_ = viper.BindPFlag(config.AssessmentURLFlag, cmd.Flags().Lookup(config.AssessmentURLFlag))
	_ = viper.BindPFlag(config.CloudServiceIDFlag, cmd.Flags().Lookup(config.CloudServiceIDFlag))
	_ = viper.BindPFlag(config.DiscoveryAutoStartFlag, cmd.Flags().Lookup(config.DiscoveryAutoStartFlag))
	_ = viper.BindPFlag(config.DiscoveryProviderFlag, cmd.Flags().Lookup(config.DiscoveryProviderFlag))
	_ = viper.BindPFlag(config.DiscoveryResourceGroupFlag, cmd.Flags().Lookup(config.DiscoveryResourceGroupFlag))
	_ = viper.BindPFlag(config.DiscoveryCSAFDomainFlag, cmd.Flags().Lookup(config.DefaultCSAFDomain))
	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.Flags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.Flags().Lookup(config.APIHTTPPortFlag))
}
