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

package assessment

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
)

// NewListAssessmentResultsCommand returns a cobra command for the `list` subcommand
func NewListAssessmentResultsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-assessment-results",
		Short: "Lists all assessment results stored in the assessment service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  assessment.AssessmentClient
				res     *assessment.ListAssessmentResultsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = assessment.NewAssessmentClient(session)

			res, err = client.ListAssessmentResults(context.Background(), &assessment.ListAssessmentResultsRequest{})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewListStatisticsCommand returns a cobra command for the `list` subcommand
func NewListStatisticsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-statistics",
		Short: "Lists statistics in the assessment service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  assessment.AssessmentClient
				res     *assessment.ListStatisticsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = assessment.NewAssessmentClient(session)

			res, err = client.ListStatistics(context.Background(), &assessment.ListStatisticsRequest{})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewAssessmentCommand returns a cobra command for `assessment` subcommands
func NewAssessmentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assessment",
		Short: "Assessment service commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListAssessmentResultsCommand(),
		NewListStatisticsCommand(),
	)
}
