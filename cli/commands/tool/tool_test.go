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
	"os"
	"testing"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(clitest.RunCLITest(m, server.WithServices(service_orchestrator.NewService())))
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
