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
	"net"
	"net/http"

	"clouditor.io/clouditor/cli"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// URLFlag is the viper flag for the server url.
	URLFlag = "url"

	// OAuth2Server is the viper flag for the OAuth 2.0 authorization server.
	OAuth2ServerFlag = "oauth2-server"

	// OAuth2ClientIDFlag is the viper flag for the OAuth 2.0 client ID.
	OAuth2ClientIDFlag = "oauth2-client-id"

	// DefaultOAuth2Server is the default OAuth 2.0 authorization server.
	DefaultOAuth2Server = "http://localhost:8080"

	// DefaultClientID is the default OAuth 2.0 client ID for the CLI.
	DefaultClientID = "cli"

	// DefaultCallbackServerAddress is the default address for the callback server.
	DefaultCallbackServerAddress = "localhost:10000"
)

// DefaultCallback is the default callback URL of the callback server.
var DefaultCallback = fmt.Sprintf("http://%s/callback", DefaultCallbackServerAddress)

// VerifierGenerator is a function that generates a new verifier.
var VerifierGenerator = oauth2.GenerateSecret

var callbackServerReady chan bool = make(chan bool)

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
				sock    net.Listener
				code    string
				config  *oauth2.Config
				authURL string
			)

			// Retrieve the URL of our authentication server
			authURL = viper.GetString(OAuth2ServerFlag)

			// Create an OAuth 2 config
			config = &oauth2.Config{
				ClientID: DefaultClientID,
				Endpoint: oauth2.Endpoint{
					AuthURL:  fmt.Sprintf("%s/authorize", authURL),
					TokenURL: fmt.Sprintf("%s/token", authURL),
				},
				RedirectURL: DefaultCallback,
			}

			srv := newCallbackServer(config)

			//go func() {
			//	exec.Command("open", authURL).Run()
			//}()

			go func() {
				sock, err = net.Listen("tcp", srv.Addr)
				if err != nil {
					fmt.Printf("Could not start web server for OAuth 2.0 authorization code flow: %v", err)
				}
				go func() {
					callbackServerReady <- true
				}()

				err = srv.Serve(sock)
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

			if session, err = cli.NewSession(args[0], config, token); err != nil {
				return fmt.Errorf("could not connect: %w", err)
			}

			if err = session.Save(); err != nil {
				return fmt.Errorf("could not save session: %w", err)
			}

			fmt.Print("\nLogin successful\n")

			return err
		},
	}

	cmd.PersistentFlags().String(OAuth2ServerFlag, DefaultOAuth2Server, "the URL of the OAuth 2.0 server")
	_ = viper.BindPFlag(OAuth2ServerFlag, cmd.PersistentFlags().Lookup(OAuth2ServerFlag))

	cmd.PersistentFlags().String(OAuth2ClientIDFlag, DefaultClientID, "the OAuth 2.0 client ID")
	_ = viper.BindPFlag(OAuth2ClientIDFlag, cmd.PersistentFlags().Lookup(OAuth2ClientIDFlag))

	return cmd
}

type callbackServer struct {
	http.Server

	verifier string
	config   *oauth2.Config
	code     chan string
}

func newCallbackServer(config *oauth2.Config) *callbackServer {
	var mux = http.NewServeMux()

	var srv = &callbackServer{
		Server: http.Server{
			Handler: mux,
			Addr:    DefaultCallbackServerAddress,
		},
		verifier: VerifierGenerator(),
		config:   config,
		code:     make(chan string),
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
