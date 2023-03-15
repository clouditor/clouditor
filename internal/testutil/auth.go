package testutil

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	oauth2 "github.com/oxisto/oauth2go"
	"github.com/oxisto/oauth2go/login"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc/metadata"
)

const (
	TestAuthUser         = "clouditor"
	TestAuthPassword     = "clouditor"
	TestCustomClaims     = "cloudserviceid"
	TestAllowAllClaims   = "cladmin"
	TestAuthClientID     = "client"
	TestAuthClientSecret = "secret"

	TestCloudService1     = "11111111-1111-1111-1111-111111111111"
	TestCloudService2     = "22222222-2222-2222-2222-222222222222"
	TestCloudServiceName1 = "My Cloud Service 1"
	TestCloudServiceName2 = "My Cloud Service 2"
)

var (
	// TestContextOnlyService1 is an incoming context with a JWT that only allows access to cloud service ID
	// 11111111-1111-1111-1111-111111111111
	TestContextOnlyService1 context.Context

	// TestContextOnlyService1 is an incoming context with a JWT that allows access to all cloud services
	TestContextAllowAll context.Context

	// TestBrokenContext contains an invalid JWT
	TestBrokenContext = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer what",
	}))

	// TestClaimsOnlyService1 contains claims that authorize the user for the cloud service
	// 11111111-1111-1111-1111-111111111111.
	TestClaimsOnlyService1 = jwt.MapClaims{
		"sub": "me",
		"cloudserviceid": []string{
			TestCloudService1,
		},
		"other": []int{1, 2},
	}

	// TestClaimsOnlyService1 contains claims that authorize the user for all cloud services.
	TestClaimsAllowAll = jwt.MapClaims{
		"sub":     "me",
		"cladmin": true,
	}
)

func init() {
	var (
		err   error
		token *jwt.Token
		t     string
	)

	// Create a new token instead of hard-coding one
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &TestClaimsOnlyService1)
	t, err = token.SignedString([]byte("mykey"))
	if err != nil {
		panic(err)
	}

	// Create a context containing our token
	TestContextOnlyService1 = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer " + t,
	}))

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, &TestClaimsAllowAll)
	t, err = token.SignedString([]byte("mykey"))
	if err != nil {
		panic(err)
	}

	// Create a context containing our token
	TestContextAllowAll = metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		"authorization": "bearer " + t,
	}))
}

// StartAuthenticationServer starts an authentication server on a random port with
// users and clients specified in the TestAuthUser and TestAuthClientID constants.
func StartAuthenticationServer() (srv *oauth2.AuthorizationServer, port uint16, err error) {
	var nl net.Listener

	// create a new socket for HTTP communication
	nl, err = net.Listen("tcp", ":0")
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
