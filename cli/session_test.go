package cli_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/persistence"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSession(t *testing.T) {
	var err error

	err = persistence.InitDB(true, "", 0)

	assert.Nil(t, err)

	var authService *service_auth.Service
	var session *cli.Session

	// create a new socket for gRPC communication
	sock, err := net.Listen("tcp", ":0")

	assert.Nil(t, err)

	defer sock.Close()

	assert.Nil(t, err, "could not listen")

	authService = &service_auth.Service{}
	authService.CreateDefaultUser("clouditor", "clouditor")

	server := grpc.NewServer()
	auth.RegisterAuthenticationServer(server, authService)

	go func() {
		// serve the gRPC socket
		err = server.Serve(sock)
		assert.Nil(t, err)
	}()

	session, err = cli.NewSession(fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port), "test")

	assert.Nil(t, err)
	assert.NotNil(t, session)

	fmt.Printf("%+v\n", session)

	client := auth.NewAuthenticationClient(session)

	var response *auth.LoginResponse

	// login with real user
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})

	assert.Nil(t, err)
	assert.NotNil(t, response)

	// login with non-existing user
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "some-other-user", Password: "password"})

	assert.NotNil(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Nil(t, response)
}
