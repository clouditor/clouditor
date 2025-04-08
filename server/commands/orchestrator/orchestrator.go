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

package orchestrator

import (
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewOrchestratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "Starts a server which contains the Clouditor Orchestrator Service",
		Long:  "This command starts a Clouditor Orchestrator service",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher(cmd.Use,
				orchestrator.DefaultServiceSpec(),
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
	cmd.Flags().Bool(config.CreateDefaultTargetOfEvaluationFlag, config.DefaultCreateDefaultTarget, "Creates a default target of evaluation if it does not exist")
	cmd.Flags().String(config.DefaultTargetOfEvaluationNameFlag, config.DefaultTargetOfEvaluationName, "Name of the default target of evaluation")
	cmd.Flags().String(config.DefaultTargetOfEvaluationDescriptionFlag, config.DefaultTargetOfEvaluationDescription, "Description of the default target of evaluation")
	cmd.Flags().Int32(config.DefaultTargetOfEvaluationTypeFlag, int32(config.DefaultTargetOfEvaluationType), "Type of the default target of evaluation; (1=cloud, 2=product, 3=organisation)")
	if cmd.Flag(config.APIgRPCPortFlag) == nil {
		cmd.Flags().Uint16(config.APIgRPCPortFlag, config.DefaultAPIgRPCPortOrchestrator, "Specifies the port used for the Clouditor gRPC API")
	}
	if cmd.Flag(config.APIHTTPPortFlag) == nil {
		cmd.Flags().Uint16(config.APIHTTPPortFlag, config.DefaultAPIHTTPPortOrchestrator, "Specifies the port used for the Clouditor HTTP API")
	}

	_ = viper.BindPFlag(config.CreateDefaultTargetOfEvaluationFlag, cmd.Flags().Lookup(config.CreateDefaultTargetOfEvaluationFlag))
	_ = viper.BindPFlag(config.DefaultTargetOfEvaluationNameFlag, cmd.Flags().Lookup(config.DefaultTargetOfEvaluationNameFlag))
	_ = viper.BindPFlag(config.DefaultTargetOfEvaluationDescriptionFlag, cmd.Flags().Lookup(config.DefaultTargetOfEvaluationDescriptionFlag))
	_ = viper.BindPFlag(config.DefaultTargetOfEvaluationTypeFlag, cmd.Flags().Lookup(config.DefaultTargetOfEvaluationTypeFlag))
	_ = viper.BindPFlag(config.APIgRPCPortFlag, cmd.Flags().Lookup(config.APIgRPCPortFlag))
	_ = viper.BindPFlag(config.APIHTTPPortFlag, cmd.Flags().Lookup(config.APIHTTPPortFlag))

}
