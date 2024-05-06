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

package assessment

import (
	"fmt"

	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/service/assessment"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAssessmentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assessment",
		Short: "Starts a server which contains the Clouditor Assessment Service",
		Long:  "This command starts a Clouditor Assessment service",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				assessment.DefaultServiceSpec(),
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
	// Set the OrchestratorURLFlag default value to the default orchestrator gRPC port, e.g., "localhost:9090"
	if cmd.Flag(config.OrchestratorURLFlag) == nil {
		cmd.Flags().String(config.OrchestratorURLFlag, config.DefaultOrchestratorURL, "Specifies the Orchestrator URL")
	}
	// Set the EvidenceStoreURLFLag default value to the default evidence store gRPC port, e.g., "localhost:9092"
	cmd.Flags().String(config.EvidenceStoreURLFlag, fmt.Sprintf("localhost:%q", config.DefaultAPIgRPCPortEvidenceStore), "Specifies the Evidence Store URL")

	if cmd.Flag(config.APIgRPCPortFlag) == nil {
		cmd.Flags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPortAssessment, "Specifies the port used for the Clouditor gRPC API")
	}
	if cmd.Flag(config.APIHTTPPortFlag) == nil {
		cmd.Flags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPortAssessment, "Specifies the port used for the Clouditor HTTP API")
	}

	_ = viper.BindPFlag(config.OrchestratorURLFlag, cmd.Flags().Lookup(config.OrchestratorURLFlag))
	_ = viper.BindPFlag(config.EvidenceStoreURLFlag, cmd.Flags().Lookup(config.EvidenceStoreURLFlag))
	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.Flags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.Flags().Lookup(config.APIHTTPPortFlag))
}
