// Copyright 2021 Fraunhofer AISEC
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

package cloud

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewCreateTargetOfEvaluationCommand returns a cobra command for the `discover` subcommand
func NewCreateTargetOfEvaluationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name]",
		Short: "Registers a new target of evaluation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.TargetOfEvaluation
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			name := args[0]

			res, err = client.CreateTargetOfEvaluation(context.Background(), &orchestrator.CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
					Name: name,
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"aws", "azure"}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewListTargetsOfEvaluationCommand returns a cobra command for the `list` subcommand
func NewListTargetsOfEvaluationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all target of evaluations",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListTargetsOfEvaluationResponse
				target  []*orchestrator.TargetOfEvaluation
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			target, err = api.ListAllPaginated(&orchestrator.ListTargetsOfEvaluationRequest{}, client.ListTargetsOfEvaluation, func(res *orchestrator.ListTargetsOfEvaluationResponse) []*orchestrator.TargetOfEvaluation {
				return res.Targets
			})

			// Build a response with all target of evaluations
			res = &orchestrator.ListTargetsOfEvaluationResponse{
				Targets: target,
			}

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewGetTargetOfEvaluationCommand returns a cobra command for the `get` subcommand
func NewGetTargetOfEvaluationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Retrieves a target of evaluation by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.TargetOfEvaluation
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]

			res, err = client.GetTargetOfEvaluation(context.Background(), &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: targetID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewRemoveTargetOfEvaluationComand returns a cobra command for the `get` subcommand
func NewRemoveTargetOfEvaluationComand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [id]",
		Short: "Removes a target of evaluation by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *emptypb.Empty
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]

			res, err = client.RemoveTargetOfEvaluation(context.Background(), &orchestrator.RemoveTargetOfEvaluationRequest{TargetOfEvaluationId: targetID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	return cmd
}

// NewUpdateTargetOfEvaluationCommand returns a cobra command for the `update` subcommand
func NewUpdateTargetOfEvaluationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a target of evaluation",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.TargetOfEvaluation
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.UpdateTargetOfEvaluation(context.Background(), &orchestrator.UpdateTargetOfEvaluationRequest{
				TargetOfEvaluation: &orchestrator.TargetOfEvaluation{
					Id:          viper.GetString("id"),
					Name:        viper.GetString("name"),
					Description: viper.GetString("description"),
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	cmd.PersistentFlags().String("id", "", "the target of evaluation id to update")
	cmd.PersistentFlags().StringP("name", "n", "", "the name of the target of evaluation")
	cmd.PersistentFlags().StringP("description", "d", "", "an optional description")

	_ = cmd.MarkPersistentFlagRequired("id")
	_ = cmd.MarkPersistentFlagRequired("name")
	_ = viper.BindPFlag("id", cmd.PersistentFlags().Lookup("id"))
	_ = viper.BindPFlag("name", cmd.PersistentFlags().Lookup("name"))
	_ = viper.BindPFlag("description", cmd.PersistentFlags().Lookup("description"))
	_ = viper.BindPFlag("control-ids", cmd.PersistentFlags().Lookup("control-ids"))

	_ = cmd.RegisterFlagCompletionFunc("id", cli.ValidArgsGetTargetOfEvaluation)
	_ = cmd.RegisterFlagCompletionFunc("name", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("description", cli.DefaultArgsShellComp)

	return cmd
}

// NewGetMetricConfigurationCommand returns a cobra command for the `get-metric-configuration` subcommand
// TODO(oxisto): Can we have something like cl cloud get <id> metric-configuration <id>?
func NewGetMetricConfigurationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-metric-configuration",
		Short: "Retrieves a metric configuration for a specific target of evaluation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *assessment.MetricConfiguration
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			targetID := args[0]
			metricID := args[1]

			res, err = client.GetMetricConfiguration(context.Background(), &orchestrator.GetMetricConfigurationRequest{TargetOfEvaluationId: targetID, MetricId: metricID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewCloudCommand returns a cobra command for `cloud` subcommands
func NewCloudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Target target of evaluations commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewCreateTargetOfEvaluationCommand(),
		NewListTargetsOfEvaluationCommand(),
		NewGetTargetOfEvaluationCommand(),
		NewUpdateTargetOfEvaluationCommand(),
		NewRemoveTargetOfEvaluationComand(),
		NewGetMetricConfigurationCommand(),
	)
}
