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

package login

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
				client  auth.AuthenticationClient
				res     *auth.LoginResponse
				req     *auth.LoginRequest
				session *cli.Session
			)

			if session, err = cli.NewSession(args[0]); err != nil {
				return fmt.Errorf("could not connect: %w", err)
			}

			if viper.GetString("username") == "" || viper.GetString("password") == "" {
				if req, err = cli.PromptForLogin(); err != nil {
					return fmt.Errorf("could not prompt for password: %w", err)
				}
			} else {
				req = &auth.LoginRequest{
					Username: viper.GetString("username"),
					Password: viper.GetString("password"),
				}
			}

			client = auth.NewAuthenticationClient(session)

			if res, err = client.Login(context.Background(), req); err != nil {
				return fmt.Errorf("could not login: %w", err)
			}

			// update the session
			session.Token = res.AccessToken

			if err = session.Save(); err != nil {
				return fmt.Errorf("could not save session: %w", err)
			}

			fmt.Print("\nLogin successful\n")

			return err
		},
	}

	cmd.PersistentFlags().StringP("username", "u", "", "the username. if not specified, a prompt will be displayed")
	_ = viper.BindPFlag("username", cmd.PersistentFlags().Lookup("username"))

	cmd.PersistentFlags().StringP("password", "p", "", "the password. if not specified, a prompt will be displayed")
	_ = viper.BindPFlag("password", cmd.PersistentFlags().Lookup("password"))

	return cmd
}
