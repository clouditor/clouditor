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

package login

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	// URLFlag is the viper flag for the server url
	URLFlag = "url"
)

// NewLoginCommand returns a cobra command for `login` subcommands
func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [login]",
		Short: "Log in to Clouditor",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				err     error
				conn    *grpc.ClientConn
				client  auth.AuthenticationClient
				res     *auth.LoginResponse
				req     *auth.LoginRequest
				session *cli.Session
			)

			session = cli.NewSession(args[0])

			if conn, err = grpc.Dial(session.URL, grpc.WithInsecure()); err != nil {
				return fmt.Errorf("could not connect: %w", err)
			}

			if req, err = cli.PromtForLogin(); err != nil {
				return fmt.Errorf("could not prompt for password: %w", err)
			}

			client = auth.NewAuthenticationClient(conn)

			if res, err = client.Login(context.Background(), req); err != nil {
				return fmt.Errorf("could not login: %w", err)
			}

			// update the session
			session.Token = res.Token
			session.Save()

			fmt.Print("\nLogin succesful\n")

			return err
		},
	}

	return cmd
}
