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

package orchestrator

import (
	"fmt"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"

	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
)

// NewListAssessmentResultsCommand returns a cobra command for the `list-assessment-results` subcommand
func NewListAssessmentResultsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-assessment-results",
		Short: "Lists all assessment results stored in the Orchestrator service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *assessment.ListAssessmentResultsResponse
				results []*assessment.AssessmentResult
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			results, err = api.ListAllPaginated(&assessment.ListAssessmentResultsRequest{}, client.ListAssessmentResults, func(res *assessment.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
				return res.Results
			})

			// Build a response with all results
			res = &assessment.ListAssessmentResultsResponse{
				Results: results,
			}

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewOrchestratorCommand returns a cobra command for `orchestrator` subcommands
func NewOrchestratorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "Orchestrator service commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListAssessmentResultsCommand(),
	)
}
