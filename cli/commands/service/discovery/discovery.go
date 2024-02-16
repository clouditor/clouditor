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

package discovery

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/cli"

	"github.com/spf13/cobra"
)

// NewStartDiscoveryCommand returns a cobra command for the `start` subcommand
func NewStartDiscoveryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the discovery",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  discovery.DiscoveryClient
				res     *discovery.StartDiscoveryResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = discovery.NewDiscoveryClient(session)

			res, err = client.Start(context.Background(), &discovery.StartDiscoveryRequest{})

			return session.HandleResponse(res, err)
		},
	}

	return cmd
}

// NewQueryDiscoveryCommand returns a cobra command for the `start` subcommand
func NewQueryDiscoveryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Queries the discovery",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  discovery.DiscoveryClient
				res     *discovery.ListResourcesResponse
				results []*discovery.Resource
				req     discovery.ListResourcesRequest
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = discovery.NewDiscoveryClient(session)

			req = discovery.ListResourcesRequest{
				PageSize:  0,
				PageToken: "",
				OrderBy:   "",
				Asc:       false,
				Filter:    &discovery.ListResourcesRequest_Filter{},
			}

			if len(args) > 0 {
				req.Filter.Type = &args[0]
			}

			results, err = api.ListAllPaginated(&discovery.ListResourcesRequest{}, client.ListResources, func(res *discovery.ListResourcesResponse) []*discovery.Resource {
				return res.Results
			})

			// Build a response with all results
			res = &discovery.ListResourcesResponse{
				Results: results,
			}

			return session.HandleResponse(res, err)
		},
	}

	return cmd
}

// NewDiscoveryCommand returns a cobra command for `discovery` subcommands
func NewDiscoveryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discovery",
		Short: "Discovery service commands",
	}

	AddCommands(cmd)

	return cmd
}

// AddCommands adds all subcommands
func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewStartDiscoveryCommand(),
		NewQueryDiscoveryCommand(),
		NewExperimentalCommand(),
	)
}
