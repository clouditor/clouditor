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

package resource

import (
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/cli"

	"github.com/spf13/cobra"
)

// NewListResourceCommand returns a cobra command for the `start` subcommand
func NewListResourcesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  evidence.EvidenceStoreClient
				res     *evidence.ListResourcesResponse
				results []*evidence.Resource
				req     evidence.ListResourcesRequest
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = evidence.NewEvidenceStoreClient(session)

			req = evidence.ListResourcesRequest{
				PageSize:  0,
				PageToken: "",
				OrderBy:   "",
				Asc:       false,
				Filter:    &evidence.ListResourcesRequest_Filter{},
			}

			if len(args) > 0 {
				req.Filter.Type = &args[0]
			}

			results, err = api.ListAllPaginated(&evidence.ListResourcesRequest{}, client.ListResources, func(res *evidence.ListResourcesResponse) []*evidence.Resource {
				return res.Results
			})

			// Build a response with all results
			res = &evidence.ListResourcesResponse{
				Results: results,
			}

			return session.HandleResponse(res, err)
		},
	}

	return cmd
}

// NewResourceCommand returns a cobra command for `resource` subcommands
func NewResourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Resource commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListResourcesCommand(),
	)
}
