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
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/cli"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var sock net.Listener
var server *grpc.Server

func TestMain(m *testing.M) {
	var err error

	sock, server, err = service_auth.StartStandaloneAuthServer(":0")

	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestSession(t *testing.T) {
	var (
		err     error
		session *cli.Session
		dir     string
	)
	defer sock.Close()
	defer server.Stop()

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
