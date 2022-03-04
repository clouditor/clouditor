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

package login

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"testing"

	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/service"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	sock   net.Listener
	server *grpc.Server
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	sock, server, err = service.StartGRPCServer("")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	var (
		err      error
		dir      string
		verifier string
		authSrv  *oauth2.AuthorizationServer
		port     int
	)

	authSrv, port, err = testutil.StartAuthenticationServer()

	defer sock.Close()
	defer server.Stop()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	viper.Set("auth-server", fmt.Sprintf("http://localhost:%d", port))
	viper.Set("session-directory", dir)

	// Issue a code that we can use in the callback
	verifier = "012345678901234567890123456789" // TODO(oxisto): random verifier
	code := authSrv.IssueCode(oauth2.GenerateCodeChallenge(verifier))

	cmd := NewLoginCommand()

	// Simulate a callback
	go func() {
		http.Get(fmt.Sprintf("http://localhost:10000/callback?code=%s", code))
	}()

	err = cmd.RunE(nil, []string{fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)})
	assert.NoError(t, err)
}
