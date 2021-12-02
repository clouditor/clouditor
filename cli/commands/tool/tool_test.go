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

package tool_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli/commands/login"
	"clouditor.io/clouditor/cli/commands/tool"
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
	)

	sock, server, err = service_auth.StartDedicatedAuthServer(":0")
	orchestrator.RegisterOrchestratorServer(server, service_orchestrator.NewService())

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

	cmd := tool.NewListToolsCommand()
	err = cmd.RunE(nil, []string{})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "method ListAssessmentTools not implemented")
}

func TestShowTool(t *testing.T) {
	var err error

	cmd := tool.NewShowToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "method GetAssessmentTool not implemented")
}

func TestUpdateTool(t *testing.T) {
	var err error

	cmd := tool.NewUpdateToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "method UpdateAssessmentTool not implemented")
}

func TestRegisterTool(t *testing.T) {
	var err error

	cmd := tool.NewRegisterToolCommand()
	err = cmd.RunE(nil, []string{})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "method RegisterAssessmentTool not implemented")
}

func TestDeregisterTool(t *testing.T) {
	var err error

	cmd := tool.NewDeregisterToolCommand()
	err = cmd.RunE(nil, []string{"1"})

	// unsupported for now
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "method DeregisterAssessmentTool not implemented")
}
