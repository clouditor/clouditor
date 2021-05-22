/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package tool

import (
	"fmt"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// NewListToolCommand returns a cobra command for the `list` subcommand
func NewListToolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				conn    *grpc.ClientConn
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListAssessmentToolsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			if conn, err = grpc.Dial(session.URL, grpc.WithInsecure()); err != nil {
				return fmt.Errorf("could not connect: %v", err)
			}

			client = orchestrator.NewOrchestratorClient(conn)

			res, err = client.ListAssessmentTools(session.Context(), &orchestrator.ListAssessmentToolsRequest{})

			session.HandleResponse(res, err)

			return err
		},
	}

	cmd.PersistentFlags().StringP("metric-id", "m", "", "only list tools for this metric")
	viper.BindPFlag("metric-id", cmd.PersistentFlags().Lookup("metric-id"))

	return cmd
}

// NewListToolCommand returns a cobra command for the `list` subcommand
func NewShowToolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Get details of a tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				conn    *grpc.ClientConn
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListAssessmentToolsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			if conn, err = grpc.Dial(session.URL, grpc.WithInsecure()); err != nil {
				return fmt.Errorf("could not connect: %v", err)
			}

			client = orchestrator.NewOrchestratorClient(conn)

			res, err = client.ListAssessmentTools(session.Context(), &orchestrator.ListAssessmentToolsRequest{})

			session.HandleResponse(res, err)

			return err
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return getTools(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.PersistentFlags().StringP("metric-id", "m", "", "only list tools for this metric")
	viper.BindPFlag("metric-id", cmd.PersistentFlags().Lookup("metric-id"))

	return cmd
}

func getTools(toComplete string) []string {
	var (
		err     error
		session *cli.Session
		conn    *grpc.ClientConn
		client  orchestrator.OrchestratorClient
		res     *orchestrator.ListAssessmentToolsResponse
	)

	if session, err = cli.ContinueSession(); err != nil {
		fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
		return nil
	}

	if conn, err = grpc.Dial(session.URL, grpc.WithInsecure()); err != nil {
		return []string{}
	}

	client = orchestrator.NewOrchestratorClient(conn)

	if res, err = client.ListAssessmentTools(session.Context(), &orchestrator.ListAssessmentToolsRequest{}); err != nil {
		return []string{}
	}

	var tools []string
	for _, v := range res.Tools {
		tools = append(tools, fmt.Sprintf("%s\t%s: %s", v.Id, v.Name, v.Description))
	}

	return tools
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
	)
}
