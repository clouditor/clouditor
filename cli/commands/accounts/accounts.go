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

package accounts

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
)

// NewListAccountsCommand returns a cobra command for the `list` subcommand
func NewListAccountsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.ListAccountsResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			res, err = client.ListAccounts(context.Background(), &orchestrator.ListAccountsRequests{})

			return session.HandleResponse(res, err)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
	}

	return cmd
}

// NewDiscoverAccountsCommand returns a cobra command for the `discover` subcommand
func NewDiscoverAccountsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discover [type]",
		Short: "Discover an account automatically",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  orchestrator.OrchestratorClient
				res     *orchestrator.CloudAccount
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = orchestrator.NewOrchestratorClient(session)

			t := args[0]

			res, err = client.CreateAccount(context.Background(), &orchestrator.CreateAccountRequest{
				Account: &orchestrator.CloudAccount{
					AutoDiscover: true,
					AccountType:  t,
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

// NewAccountsCommand returns a cobra command for `accounts` subcommands
func NewAccountsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Accounts commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListAccountsCommand(),
		NewDiscoverAccountsCommand(),
	)
}
