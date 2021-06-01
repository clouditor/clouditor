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

package login_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/persistence"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

func TestLogin(t *testing.T) {
	var (
		err error
		dir string
	)

	defer sock.Close()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.Nil(t, err)
	assert.NotEmpty(t, dir)

	viper.Set("username", "clouditor")
	viper.Set("password", "clouditor")
	viper.Set("session-directory", dir)

	cmd := login.NewLoginCommand()
	err = cmd.RunE(nil, []string{fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)})
	assert.Nil(t, err)
}
