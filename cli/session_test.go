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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var sock net.Listener

func init() {
	StartAuthServer()
}

func StartAuthServer() {
	var (
		err         error
		authService *service_auth.Service
	)

	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", ":0") // random open port
	if err != nil {
		log.Fatalf("Could not listen: %v", err)
	}

	err = persistence.InitDB(true, "", 0)

	if err != nil {
		log.Fatalf("Server exited: %v", err)
	}

	authService = &service_auth.Service{}
	authService.CreateDefaultUser("clouditor", "clouditor")

	server := grpc.NewServer()
	auth.RegisterAuthenticationServer(server, authService)

	go func() {
		// serve the gRPC socket
		_ = server.Serve(sock)
	}()
}

// TODO(oxisto): instead of replicating the things we do in the cmd, it would be good to call the command with arguments
func TestSession(t *testing.T) {
	var (
		err     error
		session *cli.Session
		dir     string
	)
	defer sock.Close()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.Nil(t, err)
	assert.NotEmpty(t, dir)

	viper.Set("session-directory", dir)

	session, err = cli.NewSession(fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer session.Close()

	assert.Nil(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, dir, session.Folder)

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

	session, err = cli.ContinueSession()
	assert.Nil(t, err)
	assert.NotNil(t, session)

	client = auth.NewAuthenticationClient(session)

	// login with non-existing user
	// TODO(oxisto): Should be moved to a service/auth test. here we should only test the session mechanism
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "some-other-user", Password: "password"})

	assert.NotNil(t, err)

	s, ok := status.FromError(err)

	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, s.Code())
	assert.Nil(t, response)
}
