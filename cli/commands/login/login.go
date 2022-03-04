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
	"net/http"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/cli"
	oauth2 "github.com/oxisto/oauth2go"
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
				session *cli.Session
				code    string
			)

			if session, err = cli.NewSession(args[0]); err != nil {
				return fmt.Errorf("could not connect: %w", err)
			}

			srv := newCallbackServer(viper.GetString("auth-server"))

			//go func() {
			//	exec.Command("open", authURL).Run()
			//}()

			go func() {
				err = srv.ListenAndServe()
				if err != http.ErrServerClosed {
					fmt.Printf("Could not start web server for OAuth 2.0 authorization code flow: %v", err)
					return
				}
			}()
			defer srv.Close()

			// waiting for our code
			code = <-srv.code
			token, err := srv.config.Exchange(context.Background(), code,
				oauth2.SetAuthURLParam("code_verifier", srv.verifier),
			)

			if err != nil {
				return err
			}

			// Update the session
			session.SetAuthorizer(api.NewOAuthAuthorizerFromConfig(srv.config, token))

			if err = session.Save(); err != nil {
				return fmt.Errorf("could not save session: %w", err)
			}

			fmt.Print("\nLogin successful\n")

			return err
		},
	}

	cmd.PersistentFlags().StringP("auth-server", "", "http://localhost:8080", "the URL to the auth server")
	_ = viper.BindPFlag("auth-server", cmd.PersistentFlags().Lookup("auth-server"))

	return cmd
}

type callbackServer struct {
	http.Server

	verifier string
	config   *oauth2.Config
	code     chan string
}

func newCallbackServer(url string) *callbackServer {
	var mux = http.NewServeMux()

	var srv = &callbackServer{
		Server: http.Server{
			Handler: mux,
			Addr:    "localhost:10000",
		},
		// TODO(oxisto): random verifier
		verifier: "012345678901234567890123456789",
		config: &oauth2.Config{
			ClientID: "cli",
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("%s/authorize", url),
				TokenURL: fmt.Sprintf("%s/token", url),
			},
			RedirectURL: "http://localhost:10000/callback",
		},
		code: make(chan string),
	}

	mux.HandleFunc("/callback", srv.handleCallback)

	challenge := oauth2.GenerateCodeChallenge(srv.verifier)
	authURL := srv.config.AuthCodeURL("",
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	fmt.Printf("Please open %s in your browser 🚀 to continue\n", authURL)

	return srv
}

func (srv *callbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	var err error

	_, err = w.Write([]byte("Success. You can close this browser tab now"))
	if err != nil {
		w.WriteHeader(500)
	}

	srv.code <- r.URL.Query().Get("code")
}
