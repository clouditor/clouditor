// Copyright 2023 Fraunhofer AISEC
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

package evidence

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/cli"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewDiscoveryCommand returns a cobra command for `experimental` subcommands
func NewExperimentalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Experimental discovery service commands",
	}

	AddExperimentalCommands(cmd)

	return cmd
}

// AddExperimentalCommands adds all experimental subcommands
func AddExperimentalCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		NewListGraphEdgesCommand(),
		NewUpdateResourceCommand(),
	)
}

// NewListEdgesCommand returns a cobra command for the `list-graph-edges` subcommand
func NewListGraphEdgesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-graph-edges",
		Short: "Lists graph edges",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  evidence.ExperimentalResourcesClient
				res     *evidence.ListGraphEdgesResponse
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = evidence.NewExperimentalResourcesClient(session)

			res, err = client.ListGraphEdges(context.Background(), &evidence.ListGraphEdgesRequest{})

			return session.HandleResponse(res, err)
		},
	}

	return cmd
}

// NewUpdateResourceCommand returns a cobra command for the `list-graph-edges` subcommand
func NewUpdateResourceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-resource",
		Short: "Updates a particular resource",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				session *cli.Session
				client  evidence.ExperimentalResourcesClient
				res     *evidence.Resource
				req     *evidence.UpdateResourceRequest
			)

			if session, err = cli.ContinueSession(); err != nil {
				fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
				return nil
			}

			client = evidence.NewExperimentalResourcesClient(session)

			req = new(evidence.UpdateResourceRequest)
			req.Resource = new(evidence.Resource)

			err = protojson.Unmarshal([]byte(args[0]), req.Resource)
			if err != nil {
				return session.HandleResponse(nil, err)
			}

			res, err = client.UpdateResource(context.Background(), req)

			return session.HandleResponse(res, err)
		},
	}

	return cmd
}
