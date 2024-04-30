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

package evidence

import (
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	evidences "clouditor.io/clouditor/v2/service/evidence"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewEvidenceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evidence",
		Short: "Starts a server which contains the Clouditor Evidence Store Service",
		Long:  "This command starts a Clouditor Evidence Store service",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				evidences.DefaultServiceSpec(),
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
	if cmd.Flag(config.APIgRPCPortFlag) == nil {
		cmd.Flags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPortEvidenceStore, "Specifies the port used for the Clouditor gRPC API")
	}
	if cmd.Flag(config.APIHTTPPortFlag) == nil {
		cmd.Flags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPortEvidenceStore, "Specifies the port used for the Clouditor HTTP API")
	}

	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.PersistentFlags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.PersistentFlags().Lookup(config.APIHTTPPortFlag))

}
