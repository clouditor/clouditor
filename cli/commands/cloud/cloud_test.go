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

package cloud

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/persistence"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	sock                net.Listener
	server              *grpc.Server
	authService         *service_auth.Service
	orchestratorService *service_orchestrator.Service
	target              *orchestrator.CloudService
	gormX               = new(persistence.GormX)
)

func TestMain(m *testing.M) {
	var (
		err error
		dir string
	)

	err = os.Chdir("../../../")
	if err != nil {
		panic(err)
	}

	err = gormX.Init(true, "", 0)
	if err != nil {
		panic(err)
	}

	orchestratorService = service_orchestrator.NewService(gormX)
	authService = service_auth.NewService(gormX)

	sock, server, err = authService.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

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

func TestNewCloudCommand(t *testing.T) {
	cmd := NewCloudCommand()

	assert.NotNil(t, cmd)
	assert.True(t, cmd.HasSubCommands())
}

func TestRegisterCloudServiceCommand(t *testing.T) {
	var err error
	var b bytes.Buffer
	var response orchestrator.CloudService

	cli.Output = &b

	cmd := NewRegisterCloudServiceCommand()
	err = cmd.RunE(nil, []string{"not_default"})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, "not_default", response.Name)
}

func TestListCloudServicesCommand(t *testing.T) {
	var (
		err      error
		b        bytes.Buffer
		response orchestrator.ListCloudServicesResponse
	)

	_, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	defer gormX.Reset()

	cli.Output = &b

	cmd := NewListCloudServicesCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Services)
}

func TestGetCloudServiceCommand(t *testing.T) {
	var (
		err      error
		b        bytes.Buffer
		response orchestrator.CloudService
	)

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	defer gormX.Reset()

	cli.Output = &b

	cmd := NewGetCloudServiceComand()
	err = cmd.RunE(nil, []string{target.Id})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, target.Id, response.Id)
}

func TestRemoveCloudServicesCommand(t *testing.T) {
	var (
		err      error
		b        bytes.Buffer
		response emptypb.Empty
	)

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	defer gormX.Reset()

	cli.Output = &b

	cmd := NewRemoveCloudServiceComand()
	err = cmd.RunE(nil, []string{target.Id})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)

	// Re-create default service
	_, err = orchestratorService.CreateDefaultTargetCloudService()

	assert.Nil(t, err)
}

func TestUpdateCloudServiceCommand(t *testing.T) {
	var (
		err      error
		b        bytes.Buffer
		response orchestrator.CloudService
	)

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	defer gormX.Reset()

	cli.Output = &b

	viper.Set("id", target.Id)
	viper.Set("name", "not_default")

	cmd := NewUpdateCloudServiceCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, target.Id, response.Id)
	assert.Equal(t, "not_default", response.Name)
}

func TestGetMetricConfiguration(t *testing.T) {
	var (
		err    error
		b      bytes.Buffer
		target *orchestrator.CloudService
	)

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	defer gormX.Reset()

	cli.Output = &b

	// create a new target service
	target, err = orchestratorService.RegisterCloudService(context.TODO(), &orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{Name: "myservice"}})

	assert.NotNil(t, target)
	assert.Nil(t, err)

	cmd := NewGetMetricConfigurationCommand()
	err = cmd.RunE(nil, []string{target.Id, "TransportEncryptionEnabled"})

	assert.Nil(t, err)
}
