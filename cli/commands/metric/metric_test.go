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

package metric

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/cli/commands/login"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestMain(m *testing.M) {
	var (
		sock                net.Listener
		server              *grpc.Server
		orchestratorService *service_orchestrator.Service
		authService         *service_auth.Service

		err error
		dir string
	)

	err = os.Chdir("../../../")
	if err != nil {
		panic(err)
	}

	orchestratorService = service_orchestrator.NewService()
	authService = service_auth.NewService(service_auth.WithApiKeySaveOnCreate(false))

	sock, server, err = authService.StartDedicatedServer(":0")
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

	if err != nil {
		panic(err)
	}

	defer sock.Close()
	defer server.Stop()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	if err != nil {
		panic(err)
	}

	viper.Set("username", "clouditor")
	viper.Set("password", "clouditor")
	viper.Set("session-directory", dir)

	cmd := login.NewLoginCommand()
	err = cmd.RunE(nil, []string{fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)})
	if err != nil {
		panic(err)
	}

	defer os.Exit(m.Run())
}

func TestListMetrics(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListMetricsCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	var response *orchestrator.ListMetricsResponse = &orchestrator.ListMetricsResponse{}

	err = protojson.Unmarshal(b.Bytes(), response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Metrics)
}

func TestGetMetric(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetMetricCommand()
	err = cmd.RunE(nil, []string{"TransportEncryptionEnabled"})

	assert.Nil(t, err)
}
