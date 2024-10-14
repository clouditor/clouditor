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

package standalone

import (
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	command_assessment "clouditor.io/clouditor/v2/server/commands/assessment"
	command_discovery "clouditor.io/clouditor/v2/server/commands/discovery"
	command_evaluation "clouditor.io/clouditor/v2/server/commands/evaluation"
	command_evidence "clouditor.io/clouditor/v2/server/commands/evidence"
	command_orchestrator "clouditor.io/clouditor/v2/server/commands/orchestrator"
	"clouditor.io/clouditor/v2/service/assessment"
	"clouditor.io/clouditor/v2/service/discovery"
	"clouditor.io/clouditor/v2/service/evaluation"
	"clouditor.io/clouditor/v2/service/evidence"
	"clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStandaloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "standalone",
		Short: "Starts a server which contains all Clouditor service",
		Long:  "This command starts all Clouditor services in standalone mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				assessment.DefaultServiceSpec(),
				discovery.DefaultServiceSpec(),
				evaluation.DefaultServiceSpec(),
				evidence.DefaultServiceSpec(),
				orchestrator.DefaultServiceSpec(),
			)
			if err != nil {
				return err
			}

			return l.Launch()
		},
	}

	// Set gRPC and HTTP port
	cmd.Flags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPort, "Specifies the port used for the Clouditor gRPC API")
	cmd.Flags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPort, "Specifies the port used for the Clouditor HTTP API")

	// Set the OrchestratorURLFlag default value to the default orchestrator URL "localhost:9090"
	if cmd.Flag(config.OrchestratorURLFlag) == nil {
		cmd.Flags().String(config.OrchestratorURLFlag, config.DefaultOrchestratorURL, "Specifies the Orchestrator URL")
	}

	// Set the AssessmentURLFlag default value to the default assessment gRPC URL "localhost:9090"
	if cmd.Flag(config.AssessmentURLFlag) == nil {
		cmd.Flags().String(config.AssessmentURLFlag, config.DefaultAssessmentURL, "Specifies the Assessment URL")
	}

	// Set the EvidenceStoreURLFLag default value to the default evidence store URL "localhost:9090"
	if cmd.Flag(config.EvidenceStoreURLFlag) == nil {
		cmd.Flags().String(config.EvidenceStoreURLFlag, config.DefaultEvidenceStoreURL, "Specifies the Evidence Store URL")
	}

	// Set the EmbeddedOAuth2ServerPublicHostFlag default value to the default embedded OAuth2 server host "http://localhost"
	if cmd.Flag(config.EmbeddedOAuth2ServerPublicHostFlag) == nil {
		cmd.Flags().String(config.EmbeddedOAuth2ServerPublicHostFlag, config.DefaultEmbeddedOAuth2ServerPublicHost, "Specifies the embedded OAuth 2.0 authorization server public host. Default is 'http://localhost'.")
	}

	// Set flag to start embedded OAuth2 server
	cmd.Flags().Bool(config.APIStartEmbeddedOAuth2ServerFlag, true, "Specifies whether the embedded OAuth 2.0 authorization server is started as part of the REST gateway. For production workloads, an external authorization server is recommended.")

	_ = viper.BindPFlag(config.APIStartEmbeddedOAuth2ServerFlag, cmd.Flags().Lookup(config.APIStartEmbeddedOAuth2ServerFlag))
	_ = viper.BindPFlag(config.EmbeddedOAuth2ServerPublicHostFlag, cmd.Flags().Lookup(config.EmbeddedOAuth2ServerPublicHostFlag))

	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.Flags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.Flags().Lookup(config.APIHTTPPortFlag))
	_ = viper.BindPFlag(config.OrchestratorURLFlag, cmd.Flags().Lookup(config.OrchestratorURLFlag))
	_ = viper.BindPFlag(config.AssessmentURLFlag, cmd.Flags().Lookup(config.AssessmentURLFlag))
	_ = viper.BindPFlag(config.EvidenceStoreURLFlag, cmd.Flags().Lookup(config.EvidenceStoreURLFlag))

	command_evidence.BindFlags(cmd)
	command_assessment.BindFlags(cmd)
	command_discovery.BindFlags(cmd)
	command_evaluation.BindFlags(cmd)
	command_orchestrator.BindFlags(cmd)

	return cmd
}
