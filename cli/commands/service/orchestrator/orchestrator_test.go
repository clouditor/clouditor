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

package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/server"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	var (
		svc *service_orchestrator.Service
		err error
	)

	clitest.AutoChdir()

	svc = service_orchestrator.NewService()

	// Store an assessment result so that output of CMD 'list' is not empty
	_, err = svc.StoreAssessmentResult(context.TODO(), &orchestrator.StoreAssessmentResultRequest{
		Result: &assessment.AssessmentResult{
			Id:                   testdata.MockTargetOfEvaluationID1,
			MetricId:             testdata.MockMetricID1,
			EvidenceId:           testdata.MockTargetOfEvaluationID1,
			TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
			Timestamp:            timestamppb.Now(),
			ResourceId:           string(testdata.MockResourceID1),
			ResourceTypes:        []string{"ResourceType"},
			ComplianceComment:    "Some comment",
			Compliant:            false,
			MetricConfiguration: &assessment.MetricConfiguration{
				TargetValue:          toStruct(1.0),
				Operator:             "<",
				IsDefault:            true,
				MetricId:             testdata.MockMetricID1,
				TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
			},
			ToolId: util.Ref(assessment.AssessmentToolId),
		}})
	if err != nil {
		panic(err)
	}

	// Store our mock catalog
	_, err = svc.CreateCatalog(context.TODO(), &orchestrator.CreateCatalogRequest{Catalog: orchestratortest.NewCatalog()})
	if err != nil {
		panic(err)
	}

	os.Exit(clitest.RunCLITest(m, server.WithServices(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewOrchestratorCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewListResultsCommand was added
	for _, v := range cmd.Commands() {
		if v.Use == "list-assessment-results" {
			return
		}
	}
	t.Errorf("No list command was added")
}

func TestNewListResultsCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListAssessmentResultsCommand()
	err := cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &orchestrator.ListAssessmentResultsResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}

func TestNewListCatalogCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListCatalogsCommand()
	err := cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &orchestrator.ListCatalogsResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Catalogs)
}

func TestNewGetCatalogCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetCatalogCommand()
	err := cmd.RunE(nil, []string{testdata.MockCatalogID1})
	assert.NoError(t, err)

	var response = &orchestrator.Catalog{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, testdata.MockCatalogID1, response.Id)
}

func TestNewGetCategoryCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetCategoryCommand()
	err := cmd.RunE(nil, []string{testdata.MockCatalogID1, testdata.MockCategoryName})
	assert.NoError(t, err)

	var response = &orchestrator.Category{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, testdata.MockCategoryName, response.Name)
}

func TestNewGetControlCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetControlCommand()
	err := cmd.RunE(nil, []string{testdata.MockCatalogID1, testdata.MockCategoryName, testdata.MockControlID1})
	assert.NoError(t, err)

	var response = &orchestrator.Control{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, testdata.MockControlID1, response.Id)
}

func toStruct(f float32) (s *structpb.Value) {
	var (
		b   []byte
		err error
	)

	s = new(structpb.Value)

	b, err = json.Marshal(f)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(b, &s); err != nil {
		return nil
	}

	return
}
