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

package tool

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/emptypb"
)

// NewListToolsCommand returns a cobra command for the `list` subcommand
func NewListToolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListAssessmentToolsResponse
				tools   []*orchestrator.AssessmentTool
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			tools, err = api.ListAllPaginated(&orchestrator.ListAssessmentToolsRequest{}, client.ListAssessmentTools, func(res *orchestrator.ListAssessmentToolsResponse) []*orchestrator.AssessmentTool {
				return res.Tools
			})

			// Build a response with all metrics
			res = &orchestrator.ListAssessmentToolsResponse{
				Tools: tools,
			}

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	cmd.PersistentFlags().StringP("metric-id", "m", "", "only list tools for this metric")
	_ = viper.BindPFlag("metric-id", cmd.PersistentFlags().Lookup("metric-id"))
	_ = cmd.RegisterFlagCompletionFunc("metric-id", cli.ValidArgsGetMetrics)

	return cmd
}

// NewShowToolCommand returns a cobra command for the `show` subcommand
func NewShowToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [id]",
		Short: "Get details of a tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.AssessmentTool
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.GetAssessmentTool(context.Background(), &orchestrator.GetAssessmentToolRequest{
				ToolId: args[0],
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetTools,
	}

	return cmd
}

// NewRegisterToolCommand returns a cobra command for the `register` subcommand
func NewRegisterToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Registers a new assessment tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.AssessmentTool
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.RegisterAssessmentTool(context.Background(), &orchestrator.RegisterAssessmentToolRequest{
				Tool: &orchestrator.AssessmentTool{
					Name:             viper.GetString("name"),
					Description:      viper.GetString("description"),
					AvailableMetrics: viper.GetStringSlice("metric-ids"),
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	cmd.PersistentFlags().StringP("name", "n", "", "the name of the tool")
	cmd.PersistentFlags().StringP("description", "d", "", "an optional description")
	cmd.PersistentFlags().StringSliceP("metric-ids", "m", []string{}, "the metric this tool assesses")
	_ = cmd.MarkPersistentFlagRequired("name")
	_ = cmd.MarkPersistentFlagRequired("metric-ids")
	_ = viper.BindPFlag("name", cmd.PersistentFlags().Lookup("name"))
	_ = viper.BindPFlag("description", cmd.PersistentFlags().Lookup("description"))
	_ = viper.BindPFlag("metric-ids", cmd.PersistentFlags().Lookup("metric-ids"))

	_ = cmd.RegisterFlagCompletionFunc("name", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("description", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("metric-ids", cli.ValidArgsGetMetrics)

	return cmd
}

// NewUpdateToolCommand returns a cobra command for the `update` subcommand
func NewUpdateToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [id]",
		Short: "Updates a registered assessment tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.AssessmentTool
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.UpdateAssessmentTool(context.Background(), &orchestrator.UpdateAssessmentToolRequest{
				Tool: &orchestrator.AssessmentTool{
					Id:               args[0],
					Name:             viper.GetString("name"),
					Description:      viper.GetString("description"),
					AvailableMetrics: viper.GetStringSlice("metric-ids"),
				},
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.DefaultArgsShellComp,
	}

	cmd.PersistentFlags().StringP("name", "n", "", "the name of the tool")
	cmd.PersistentFlags().StringP("description", "d", "", "an optional description")
	cmd.PersistentFlags().StringSliceP("metric-ids", "m", []string{}, "the metric this tool assesses")
	_ = cmd.MarkPersistentFlagRequired("name")
	_ = cmd.MarkPersistentFlagRequired("metric-ids")
	_ = viper.BindPFlag("name", cmd.PersistentFlags().Lookup("name"))
	_ = viper.BindPFlag("description", cmd.PersistentFlags().Lookup("description"))
	_ = viper.BindPFlag("metric-ids", cmd.PersistentFlags().Lookup("metric-ids"))

	_ = cmd.RegisterFlagCompletionFunc("name", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("description", cli.DefaultArgsShellComp)
	_ = cmd.RegisterFlagCompletionFunc("metric-ids", cli.ValidArgsGetMetrics)

	return cmd
}

// NewDeregisterToolCommand returns a cobra command for the `deregister` subcommand
func NewDeregisterToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reregister [id]",
		Short: "Deregisteres a tool",
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

			res, err = client.DeregisterAssessmentTool(context.Background(), &orchestrator.DeregisterAssessmentToolRequest{
				ToolId: args[0],
			})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetTools,
	}

	return cmd
}

// NewToolCommand returns a cobra command for `tool` subcommands
func NewToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tool",
		Short: "Tool commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListToolsCommand(),
		NewShowToolCommand(),
		NewRegisterToolCommand(),
	)
}
