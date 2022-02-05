// Copyright 2022 Fraunhofer AISEC
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

package service

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
