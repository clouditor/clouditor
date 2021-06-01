package cli_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
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

// TODO(oxisto): instead of replicating the things we do in the cmd, it would be good to call the command with arguments
func TestSession(t *testing.T) {
	var (
		err     error
		session *cli.Session
		dir     string
	)

	assert.Nil(t, err)

	defer sock.Close()

	assert.Nil(t, err, "could not listen")

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.Nil(t, err)
	assert.NotEmpty(t, dir)

	session, err = cli.NewSession("bufnet", dir, grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer session.Close()

	assert.Nil(t, err)
	assert.NotNil(t, session)

	fmt.Printf("%+v\n", session)

	client := auth.NewAuthenticationClient(session)

	var response *auth.LoginResponse

	// login with real user
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)

	// update the session
	session.Token = response.Token

	err = session.Save()

	assert.Nil(t, err)

	// TODO(oxisto): not quite sure how to test continue session with the bufnet dialer
	// session, err = cli.ContinueSession(dir)
	// assert.Nil(t, err)
	// assert.NotNil(t, session)

	// client = auth.NewAuthenticationClient(session)

	// login with non-existing user
	// TODO(oxisto): Should be moved to a service/auth test. here we should only test the session mechanism
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "some-other-user", Password: "password"})

	assert.NotNil(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Nil(t, response)
}
