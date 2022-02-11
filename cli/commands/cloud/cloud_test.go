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
	"clouditor.io/clouditor/service"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/cli/commands/login"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../../../")
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
	var (
		response orchestrator.CloudService

		err error
		b   bytes.Buffer
	)
	_, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

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
		response orchestrator.ListCloudServicesResponse

		err error
		b   bytes.Buffer
	)
	orchestratorService, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

	_, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)

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
		response orchestrator.CloudService
		target   *orchestrator.CloudService

		err error
		b   bytes.Buffer
	)
	orchestratorService, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	fmt.Println("target:", target)
	// target should be non-nil since it has been newly created
	assert.NotNil(t, target)
	assert.Nil(t, err)

	cli.Output = &b

	cmd := NewGetCloudServiceCommand()
	err = cmd.RunE(nil, []string{target.Id})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, target.Id, response.Id)
}

func TestRemoveCloudServicesCommand(t *testing.T) {
	var (
		response emptypb.Empty
		target   *orchestrator.CloudService

		err error
		b   bytes.Buffer
	)
	orchestratorService, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)

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
		response orchestrator.CloudService
		target   *orchestrator.CloudService

		err error
		b   bytes.Buffer
	)
	const (
		notDefault = "not_default"
	)
	orchestratorService, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)

	cli.Output = &b

	viper.Set("id", target.Id)
	viper.Set("name", notDefault)
	//viper.Set("description", "newD")

	cmd := NewUpdateCloudServiceCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, target.Id, response.Id)
	assert.Equal(t, notDefault, response.Name)
}

func TestGetMetricConfiguration(t *testing.T) {
	var (
		target *orchestrator.CloudService

		err error
		b   bytes.Buffer
	)
	orchestratorService, server, sock := startServer()
	defer sock.Close()
	defer server.Stop()

	target, err = orchestratorService.CreateDefaultTargetCloudService()
	assert.Nil(t, err)
	// target should be not nil since there are no stored cloud services yet
	assert.NotNil(t, target)

	cli.Output = &b

	// create a new target service
	target, err = orchestratorService.RegisterCloudService(context.TODO(), &orchestrator.RegisterCloudServiceRequest{Service: &orchestrator.CloudService{Name: "myservice"}})

	assert.NotNil(t, target)
	assert.Nil(t, err)

	cmd := NewGetMetricConfigurationCommand()
	err = cmd.RunE(nil, []string{target.Id, "TransportEncryptionEnabled"})

	assert.Nil(t, err)
}

// startServer starts server with orchestrator and dedicated auth server. We don't do it in TestMain since you
// can only register a service - once before server.serve(). And we do need to add new Orchestrator service because
// the DB won't be reset otherwise.
func startServer() (orchestratorService *service_orchestrator.Service, server *grpc.Server, sock net.Listener) {
	var (
		err error
		dir string
	)

	orchestratorService = service_orchestrator.NewService()

	sock, server, _, err = service.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}
	orchestrator.RegisterOrchestratorServer(server, orchestratorService)

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

	return
}
