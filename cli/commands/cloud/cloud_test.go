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
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/service"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

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
		svc      *service_orchestrator.Service
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		srv      *grpc.Server

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	cli.Output = &b

	cmd := NewRegisterCloudServiceCommand()
	err = cmd.RunE(nil, []string{"not_default"})

	assert.NoError(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, "not_default", response.Name)
}

func TestListCloudServicesCommand(t *testing.T) {
	var (
		response orchestrator.ListCloudServicesResponse
		svc      *service_orchestrator.Service
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		srv      *grpc.Server

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	_, err = svc.CreateDefaultTargetCloudService()
	assert.NoError(t, err)

	cli.Output = &b

	cmd := NewListCloudServicesCommand()
	err = cmd.RunE(nil, []string{})

	assert.NoError(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Services)
}

func TestGetCloudServiceCommand(t *testing.T) {
	var (
		response orchestrator.CloudService
		target   *orchestrator.CloudService
		svc      *service_orchestrator.Service
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		srv      *grpc.Server

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	target, err = svc.CreateDefaultTargetCloudService()

	fmt.Println("target:", target)
	// target should be non-nil since it has been newly created
	assert.NotNil(t, target)
	assert.NoError(t, err)

	cli.Output = &b

	cmd := NewGetCloudServiceCommand()
	err = cmd.RunE(nil, []string{target.Id})

	assert.NoError(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, target.Id, response.Id)
}

func TestRemoveCloudServicesCommand(t *testing.T) {
	var (
		response emptypb.Empty
		target   *orchestrator.CloudService
		svc      *service_orchestrator.Service
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		srv      *grpc.Server

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	target, err = svc.CreateDefaultTargetCloudService()
	assert.NoError(t, err)

	cli.Output = &b

	cmd := NewRemoveCloudServiceComand()
	err = cmd.RunE(nil, []string{target.Id})

	assert.NoError(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.NoError(t, err)

	// Re-create default service
	_, err = svc.CreateDefaultTargetCloudService()

	assert.NoError(t, err)
}

func TestUpdateCloudServiceCommand(t *testing.T) {
	var (
		response orchestrator.CloudService
		target   *orchestrator.CloudService
		svc      *service_orchestrator.Service
		tmpDir   string
		auth     *oauth2.AuthorizationServer
		srv      *grpc.Server

		err error
		b   bytes.Buffer
	)

	const (
		notDefault = "not_default"
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	target, err = svc.CreateDefaultTargetCloudService()
	assert.NoError(t, err)

	cli.Output = &b

	viper.Set("id", target.Id)
	viper.Set("name", notDefault)
	//viper.Set("description", "newD")

	cmd := NewUpdateCloudServiceCommand()
	err = cmd.RunE(nil, []string{})

	assert.NoError(t, err)

	err = protojson.Unmarshal(b.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, target.Id, response.Id)
	assert.Equal(t, notDefault, response.Name)
}

func TestGetMetricConfiguration(t *testing.T) {
	var (
		target *orchestrator.CloudService
		svc    *service_orchestrator.Service
		tmpDir string
		auth   *oauth2.AuthorizationServer
		srv    *grpc.Server

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	tmpDir, auth, srv = startServer(service.WithOrchestrator(svc))
	defer cleanup(tmpDir, srv, auth)

	target, err = svc.CreateDefaultTargetCloudService()
	assert.NoError(t, err)
	// target should be not nil since there are no stored cloud services yet
	assert.NotNil(t, target)

	cli.Output = &b

	// create a new target service
	target, err = svc.RegisterCloudService(context.TODO(), &orchestrator.RegisterCloudServiceRequest{CloudService: &orchestrator.CloudService{Name: "myservice"}})

	assert.NotNil(t, target)
	assert.NoError(t, err)

	cmd := NewGetMetricConfigurationCommand()
	err = cmd.RunE(nil, []string{target.Id, "TransportEncryptionEnabled"})

	assert.NoError(t, err)
}

// startServer starts a gRPC server with an orchestrator and auth service. We don't do it in TestMain since you
// can only register a service - once before server.serve(). And we do need to add new Orchestrator service because
// the DB won't be reset otherwise.
func startServer(opts ...service.StartGRPCServerOption) (tmpDir string, auth *oauth2.AuthorizationServer, srv *grpc.Server) {
	var (
		err      error
		grpcPort uint16
		authPort uint16
		sock     net.Listener
	)

	auth, authPort, err = testutil.StartAuthenticationServer()
	if err != nil {
		panic(err)
	}

	sock, srv, err = service.StartGRPCServer(testutil.JWKSURL(authPort), opts...)
	if err != nil {
		panic(err)
	}

	grpcPort = sock.Addr().(*net.TCPAddr).AddrPort().Port()

	tmpDir, err = clitest.PrepareSession(authPort, auth, fmt.Sprintf("localhost:%d", grpcPort))
	if err != nil {
		panic(err)
	}

	return
}

func cleanup(tmpDir string, srv *grpc.Server, auth *oauth2.AuthorizationServer) {
	// Remove temporary session directory
	os.RemoveAll(tmpDir)

	srv.Stop()

	auth.Close()
}
