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

package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"clouditor.io/clouditor/persistence"
	service_auth "clouditor.io/clouditor/service/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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

	// Start at least an orchestrator server, so that we have something to forward
	sock, server, err = service_auth.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	go RunServer(
		context.Background(),
		sock.Addr().(*net.TCPAddr).Port,
		0,
		WithAllowedOrigins([]string{"clouditor.io"}),
	)

	defer sock.Close()
	defer server.Stop()
	defer os.Exit(m.Run())
}

func TestCORS(t *testing.T) {
	// wait until server is ready to serve
	select {
	case <-ready:
		break
	case <-time.After(10 * time.Second):
		log.Println("Timeout while waiting for REST API")
	}

	assert.NotNil(t, sock)

	port, err := GetServerPort()

	assert.Nil(t, err)
	assert.NotEqual(t, 0, port)

	client := &http.Client{}

	req, err := http.NewRequest("OPTIONS", fmt.Sprintf("http://localhost:%d/v1/auth/login", port), nil)
	assert.Nil(t, err)
	assert.NotNil(t, req)

	req.Header.Add("Origin", "clouditor.io")
	req.Header.Add("Access-Control-Request-Method", "POST")
	resp, err := client.Do(req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, "clouditor.io", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, strings.Join(DefaultAllowedHeaders, ","), resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, strings.Join(DefaultAllowedMethods, ","), resp.Header.Get("Access-Control-Allow-Methods"))
}
