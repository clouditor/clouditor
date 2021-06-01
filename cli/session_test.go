package cli_test

import (
	"context"
	"fmt"
	"log"
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
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var sock *bufconn.Listener

func init() {
	// create a new socket for gRPC communication
	sock = bufconn.Listen(bufSize)

	var err error
	err = persistence.InitDB(true, "", 0)

	if err != nil {
		log.Fatalf("Server exited: %v", err)
	}

	var authService *service_auth.Service

	authService = &service_auth.Service{}
	authService.CreateDefaultUser("clouditor", "clouditor")

	server := grpc.NewServer()
	auth.RegisterAuthenticationServer(server, authService)

	go func() {
		// serve the gRPC socket
		_ = server.Serve(sock)
		/*if err != nil {
			log.Fatalf("Server exited: %v", err)
		}*/
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return sock.Dial()
}

func TestSession(t *testing.T) {
	var err error

	var session *cli.Session

	assert.Nil(t, err)

	defer sock.Close()

	assert.Nil(t, err, "could not listen")

	//session, err = cli.NewSession(fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port), "test")

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	assert.Nil(t, err)
	//assert.NotNil(t, session)

	fmt.Printf("%+v\n", session)

	client := auth.NewAuthenticationClient(conn)

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
