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

package metric

import (
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// NewListMetricCommand returns a cobra command for the `list` subcommand
func NewListMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListMetricsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.ListMetrics(context.Background(), &orchestrator.ListMetricsRequest{})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewGetMetricCommand returns a cobra command for the `get` subcommand
func NewGetMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get metric by metricID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.GetMetricResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}


			client = orchestrator.NewOrchestratorClient(session)

			metricID, err := strconv.Atoi(args[0])
			metricIDInt32 := int32(metricID)

			if err != nil {
				fmt.Printf("Error while parsing string to int32")
				return nil
			}

			res, err = client.GetMetric(context.Background(), &orchestrator.GetMetricsRequest{MetricId: metricIDInt32})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewMetricCommand returns a cobra command for `metric` subcommands
func NewMetricCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metric",
		Short: "Metric commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListMetricsCommand(),
		NewGetMetricsCommand(),
	)
}
