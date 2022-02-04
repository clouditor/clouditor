package common

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/rest"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	grpcPort    int
	authService *service_auth.Service
)

func TestMain(m *testing.M) {
	var (
		err    error
		server *grpc.Server
		sock   net.Listener
	)

	// A small embedded DB is needed for the server
	err = persistence.InitDB(true, "", 0)
	if err != nil {
		panic(err)
	}

	// Start at least an authentication server, so that we have something to forward
	sock, server, authService, err = service_auth.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).Port

	exit := m.Run()

	sock.Close()
	server.Stop()

	os.Exit(exit)
}

func Test(t *testing.T) {
	go func() {
		err := rest.RunServer(
			context.Background(),
			grpcPort,
			0,
			rest.WithJwks(authService.GetPublicKey()),
		)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	defer rest.StopServer(context.Background())

	// meh
	time.Sleep(time.Millisecond * 300)

	port, err := rest.GetServerPort()
	assert.ErrorIs(t, err, nil)

	loginResponse, err := authService.Login(context.TODO(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})
	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, loginResponse)

	config := ConfigureAuth(WithJwksUrl(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)))
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", loginResponse.GetToken())}})

	newCtx, err := config.AuthFunc(ctx)

	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, newCtx)

	key, ok := config.Jwks.ReadOnlyKeys()["1"]

	assert.True(t, ok)
	assert.NotNil(t, key)
}
