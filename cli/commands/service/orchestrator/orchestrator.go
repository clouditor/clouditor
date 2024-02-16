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
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"

	"clouditor.io/clouditor/v2/cli"
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
				res     *orchestrator.ListAssessmentResultsResponse
				results []*assessment.AssessmentResult
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			results, err = api.ListAllPaginated(&orchestrator.ListAssessmentResultsRequest{}, client.ListAssessmentResults, func(res *orchestrator.ListAssessmentResultsResponse) []*assessment.AssessmentResult {
				return res.Results
			})

			// Build a response with all results
			res = &orchestrator.ListAssessmentResultsResponse{
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

// NewListCatalogsCommand returns a cobra command for the `list-requirements` subcommand
func NewListCatalogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-catalogs",
		Short: "Lists all catalogs",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err      error
				session  *cli.Session
				client   orchestrator.OrchestratorClient
				res      *orchestrator.ListCatalogsResponse
				catalogs []*orchestrator.Catalog
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			catalogs, err = api.ListAllPaginated(&orchestrator.ListCatalogsRequest{}, client.ListCatalogs, func(res *orchestrator.ListCatalogsResponse) []*orchestrator.Catalog {
				return res.Catalogs
			})

			// Build a response with all results
			res = &orchestrator.ListCatalogsResponse{
				Catalogs: catalogs,
			}

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewGetCatalogCommand returns a cobra command for the `get-catalog` subcommand
func NewGetCatalogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-catalog [catalog ID]",
		Short: "Retrieves a catalog by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.Catalog
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			catalogID := args[0]

			res, err = client.GetCatalog(context.Background(), &orchestrator.GetCatalogRequest{CatalogId: catalogID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetCatalogs,
	}

	return cmd
}

// NewGetCategoryCommand returns a cobra command for the `get-category` subcommand
func NewGetCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-category [catalog ID] [category name]",
		Short: "Retrieves a category by name and catalog ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.Category
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			catalogID := args[0]
			categoryName := args[1]

			res, err = client.GetCategory(context.Background(), &orchestrator.GetCategoryRequest{CatalogId: catalogID, CategoryName: categoryName})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetCategory,
	}

	return cmd
}

// NewGetControlCommand returns a cobra command for the `get-control` subcommand
func NewGetControlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-control [catalog ID] [category name] [short name]",
		Short: "Retrieves a control by its short name, its category and catalog ID",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.Control
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			catalogID := args[0]
			categoryName := args[1]
			controlID := args[2]

			res, err = client.GetControl(context.Background(), &orchestrator.GetControlRequest{CatalogId: catalogID, CategoryName: categoryName, ControlId: controlID})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: cli.ValidArgsGetControls,
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
		NewListCatalogsCommand(),
		NewGetCatalogCommand(),
		NewGetCategoryCommand(),
		NewGetControlCommand(),
	)
}
