package testutil

import (
	"fmt"
	"net"
	"net/http"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/login"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	TestAuthUser         = "clouditor"
	TestAuthPassword     = "clouditor"
	TestAuthClientID     = "client"
	TestAuthClientSecret = "secret"
)

// StartAuthenticationServer starts an authentication server on a random port with
// users and clients specified in the TestAuthUser and TestAuthClientID constants.
func StartAuthenticationServer() (srv *oauth2.AuthorizationServer, port uint16, err error) {
	var nl net.Listener

	// create a new socket for HTTP communication
	nl, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, 0, fmt.Errorf("could not listen: %w", err)
	}

	port = nl.Addr().(*net.TCPAddr).AddrPort().Port()

	srv = oauth2.NewServer(fmt.Sprintf(":%d", port),
		oauth2.WithClient("cli", "", "http://localhost:10000/callback"),
		oauth2.WithClient(TestAuthClientID, TestAuthClientSecret, ""),
		oauth2.WithPublicURL(fmt.Sprintf("http://localhost:%d", port)),
		login.WithLoginPage(
			login.WithUser(TestAuthUser, TestAuthPassword),
			login.WithBaseURL("/v1/auth"),
		),
	)

	// simulate the /v1/auth endpoints
	srv.Handler.(*http.ServeMux).Handle("/v1/auth/token", http.StripPrefix("/v1/auth", srv.Handler))
	srv.Handler.(*http.ServeMux).Handle("/v1/auth/certs", http.StripPrefix("/v1/auth", srv.Handler))

	go func() {
		_ = srv.Serve(nl)
	}()

	return srv, port, nil
}

func JWKSURL(port uint16) string {
	return fmt.Sprintf("http://localhost:%d/v1/auth/certs", port)
}

func TokenURL(port uint16) string {
	return fmt.Sprintf("http://localhost:%d/v1/auth/token", port)
}

func AuthURL(port uint16) string {
	return fmt.Sprintf("http://localhost:%d/v1/auth/authorize", port)
}

func AuthClientConfig(port uint16) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     TestAuthClientID,
		ClientSecret: TestAuthClientSecret,
		TokenURL:     TokenURL(port),
	}
}
