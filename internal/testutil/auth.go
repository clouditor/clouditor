package testutil

import (
	"context"
	"fmt"
	"net"
	"net/http"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/login"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/metadata"
)

const (
	TestAuthUser     = "clouditor"
	TestAuthPassword = "clouditor"

	TestAuthClientID     = "client"
	TestAuthClientSecret = "secret"

	TestCustomClaims  = "cloudserviceid"
	TestCloudService1 = "11111111-1111-1111-1111-111111111111"
	TestCloudService2 = "22222222-2222-2222-2222-222222222222"
)

// TestContextOnlyService1 is an incoming context with a JWT that only allows access to cloud service ID
// 11111111-1111-1111-1111-111111111111
var TestContextOnlyService1 = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
	"authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJjbG91ZHNlcnZpY2VpZCI6WyIxMTExMTExMS0xMTExLTExMTEtMTExMS0xMTExMTExMTExMTEiXSwib3RoZXIiOlsxLDJdfQ.A4_2-yRcoPui-udHifvQxB6SJj7fR1EPjBnFs0Nl80k",
}))

// TestBrokenContext contains an invalid JWT
var TestBrokenContext = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
	"authorization": "bearer what",
}))

// StartAuthenticationServer starts an authentication server on a random port with
// users and clients specified in the TestAuthUser and TestAuthClientID constants.
func StartAuthenticationServer() (srv *oauth2.AuthorizationServer, port uint16, err error) {
	var nl net.Listener

	srv = oauth2.NewServer(":0",
		oauth2.WithClient("cli", "", "http://localhost:10000/callback"),
		oauth2.WithClient(TestAuthClientID, TestAuthClientSecret, ""),
		login.WithLoginPage(
			login.WithUser(TestAuthUser, TestAuthPassword),
			login.WithBaseURL("/v1/auth"),
		),
	)

	// simulate the /v1/auth endpoints
	srv.Handler.(*http.ServeMux).Handle("/v1/auth/token", http.StripPrefix("/v1/auth", srv.Handler))

	// create a new socket for HTTP communication
	nl, err = net.Listen("tcp", srv.Addr)
	if err != nil {
		return nil, 0, fmt.Errorf("could not listen: %w", err)
	}

	go func() {
		_ = srv.Serve(nl)
	}()

	port = nl.Addr().(*net.TCPAddr).AddrPort().Port()

	return srv, port, nil
}

func JWKSURL(port uint16) string {
	return fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)
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
