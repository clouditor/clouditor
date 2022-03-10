package testutil

import (
	"fmt"
	"net"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/login"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	TestAuthUser     = "clouditor"
	TestAuthPassword = "clouditor"

	TestAuthClientID     = "clouditor"
	TestAuthClientSecret = "clouditor"
)

// StartAuthenticationServer starts an authentication server on a random port with
// users and clients specified in the TestAuthUser and TestAuthClientID constants.
func StartAuthenticationServer() (srv *oauth2.AuthorizationServer, port int, err error) {
	var nl net.Listener

	srv = oauth2.NewServer(":0",
		oauth2.WithClient("cli", "", "http://localhost:10000/callback"),
		login.WithLoginPage(login.WithUser(TestAuthUser, TestAuthPassword)),
	)

	// create a new socket for gRPC communication
	nl, err = net.Listen("tcp", srv.Addr)
	if err != nil {
		return nil, 0, fmt.Errorf("could not listen: %w", err)
	}

	go func() {
		_ = srv.Serve(nl)
	}()

	port = nl.Addr().(*net.TCPAddr).Port

	return srv, port, nil
}

func JWKSURL(port int) string {
	return fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)
}

func TokenURL(port int) string {
	return fmt.Sprintf("http://localhost:%d/v1/auth/token", port)
}

func AuthURL(port int) string {
	return fmt.Sprintf("http://localhost:%d/v1/auth/authorize", port)
}

func AuthClientConfig(port int) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     TestAuthClientID,
		ClientSecret: TestAuthClientSecret,
		TokenURL:     TokenURL(port),
	}
}
