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
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	svc *service_orchestrator.Service
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	svc = service_orchestrator.NewService()

	os.Exit(clitest.RunCLITest(m, server.WithServices(svc)))
}

func TestNewCloudCommand(t *testing.T) {
	cmd := NewCloudCommand()

	assert.NotNil(t, cmd)
	assert.True(t, cmd.HasSubCommands())
}

func TestCreateTargetOfEvaluationCommand(t *testing.T) {
	var (
		response orchestrator.TargetOfEvaluation
		svc      *service_orchestrator.Service
		err      error
		b        bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		cli.Output = &b

		cmd := NewCreateTargetOfEvaluationCommand()
		err = cmd.RunE(nil, []string{"not_default"})

		assert.NoError(t, err)

		err = protojson.Unmarshal(b.Bytes(), &response)

		assert.NoError(t, err)
		return assert.Equal(t, "not_default", response.Name)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}

func TestListTargetOfEvaluationsCommand(t *testing.T) {
	var (
		response orchestrator.ListTargetOfEvaluationsResponse
		svc      *service_orchestrator.Service

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		_, err = svc.CreateDefaultTargetOfEvaluation()
		assert.NoError(t, err)

		cli.Output = &b

		cmd := NewListTargetOfEvaluationsCommand()
		err = cmd.RunE(nil, []string{})

		assert.NoError(t, err)

		err = protojson.Unmarshal(b.Bytes(), &response)

		assert.NoError(t, err)
		return assert.NotEmpty(t, response.Targets)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}

func TestGetTargetOfEvaluationCommand(t *testing.T) {
	var (
		response orchestrator.TargetOfEvaluation
		target   *orchestrator.TargetOfEvaluation
		svc      *service_orchestrator.Service

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		target, err = svc.CreateDefaultTargetOfEvaluation()

		fmt.Println("target:", target)
		// target should be non-nil since it has been newly created
		assert.NotNil(t, target)
		assert.NoError(t, err)

		cli.Output = &b

		cmd := NewGetTargetOfEvaluationCommand()
		err = cmd.RunE(nil, []string{target.Id})

		assert.NoError(t, err)

		err = protojson.Unmarshal(b.Bytes(), &response)

		assert.NoError(t, err)
		return assert.Equal(t, target.Id, response.Id)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}

func TestRemoveTargetOfEvaluationsCommand(t *testing.T) {
	var (
		response emptypb.Empty
		target   *orchestrator.TargetOfEvaluation
		svc      *service_orchestrator.Service

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		target, err = svc.CreateDefaultTargetOfEvaluation()
		assert.NoError(t, err)

		cli.Output = &b

		cmd := NewRemoveTargetOfEvaluationComand()
		err = cmd.RunE(nil, []string{target.Id})

		assert.NoError(t, err)

		err = protojson.Unmarshal(b.Bytes(), &response)

		assert.NoError(t, err)

		// Re-create default service
		_, err = svc.CreateDefaultTargetOfEvaluation()

		return assert.NoError(t, err)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}

func TestUpdateTargetOfEvaluationCommand(t *testing.T) {
	var (
		response orchestrator.TargetOfEvaluation
		target   *orchestrator.TargetOfEvaluation
		svc      *service_orchestrator.Service

		err error
		b   bytes.Buffer
	)

	const (
		notDefault = "not_default"
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		target, err = svc.CreateDefaultTargetOfEvaluation()
		assert.NoError(t, err)

		cli.Output = &b

		viper.Set("id", target.Id)
		viper.Set("name", notDefault)

		cmd := NewUpdateTargetOfEvaluationCommand()
		err = cmd.RunE(nil, []string{})

		assert.NoError(t, err)

		err = protojson.Unmarshal(b.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, target.Id, response.Id)
		return assert.Equal(t, notDefault, response.Name)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}

func TestGetMetricConfiguration(t *testing.T) {
	var (
		target *orchestrator.TargetOfEvaluation
		svc    *service_orchestrator.Service

		err error
		b   bytes.Buffer
	)

	svc = service_orchestrator.NewService()
	_, err = clitest.RunCLITestFunc(func() bool {
		target, err = svc.CreateDefaultTargetOfEvaluation()
		assert.NoError(t, err)
		// target should be not nil since there are no stored target of evaluations yet
		assert.NotNil(t, target)

		cli.Output = &b

		// create a new target service
		target, err = svc.CreateTargetOfEvaluation(context.TODO(), &orchestrator.CreateTargetOfEvaluationRequest{TargetOfEvaluation: &orchestrator.TargetOfEvaluation{Name: "myTarget"}})

		assert.NotNil(t, target)
		assert.NoError(t, err)

		cmd := NewGetMetricConfigurationCommand()
		err = cmd.RunE(nil, []string{target.Id, "TransportEncryptionEnabled"})

		return assert.NoError(t, err)
	}, server.WithServices(svc))
	assert.NoError(t, err)
}
