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

package tool

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var sock net.Listener
var server *grpc.Server

func TestMain(m *testing.M) {
	var (
		err error
		dir string
		s   *service_orchestrator.Service
	)

	err = os.Chdir("../../../")
	if err != nil {
		panic(err)
	}

	err = persistence.InitDB(true, "", 0)
	if err != nil {
		panic(err)
	}

	s = service_orchestrator.NewService()

	sock, server, _, err = service.StartDedicatedAuthServer(":0", service_auth.WithApiKeySaveOnCreate(false))
	orchestrator.RegisterOrchestratorServer(server, s)

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

	os.Exit(m.Run())
}

func TestListTool(t *testing.T) {
	var err error

	cmd := NewListToolsCommand()
	err = cmd.RunE(nil, []string{})

	// unsupported for now
	assert.Error(t, err)
	assert.Equal(t, "method ListAssessmentTools not implemented", err.Error())
}

func TestShowTool(t *testing.T) {
	var err error

	cmd := NewShowToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.Error(t, err)
	assert.Equal(t, "method GetAssessmentTool not implemented", err.Error())
}

func TestUpdateTool(t *testing.T) {
	var err error

	cmd := NewUpdateToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.Error(t, err)
	assert.Equal(t, "method UpdateAssessmentTool not implemented", err.Error())
}

func TestRegisterTool(t *testing.T) {
	var err error

	cmd := NewRegisterToolCommand()
	err = cmd.RunE(nil, []string{})

	// unsupported for now
	assert.Error(t, err)
	assert.Equal(t, "method RegisterAssessmentTool not implemented", err.Error())
}

func TestDeregisterTool(t *testing.T) {
	var err error

	cmd := NewDeregisterToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "method DeregisterAssessmentTool not implemented")
}
