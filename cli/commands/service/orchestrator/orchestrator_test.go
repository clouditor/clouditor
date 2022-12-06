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

	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/service"

	assessmentv1 "clouditor.io/clouditor/api/assessment/v1"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/cli"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/stretchr/testify/assert"
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
		Result: &assessmentv1.AssessmentResult{
			Id:            "11111111-1111-1111-1111-111111111111",
			MetricId:      "assessmentResultMetricID",
			EvidenceId:    "11111111-1111-1111-1111-111111111111",
			Timestamp:     timestamppb.Now(),
			ResourceId:    "myResource",
			ResourceTypes: []string{"ResourceType"},
			MetricConfiguration: &assessmentv1.MetricConfiguration{
				TargetValue: toStruct(1.0),
				Operator:    "operator",
				IsDefault:   true,
			}}})
	if err != nil {
		panic(err)
	}

	// Store our mock catalog
	_, err = svc.CreateCatalog(context.TODO(), &orchestrator.CreateCatalogRequest{Catalog: orchestratortest.NewCatalog()})
	if err != nil {
		panic(err)
	}

	os.Exit(clitest.RunCLITest(m, service.WithOrchestrator(svc)))
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

	var response = &assessmentv1.ListAssessmentResultsResponse{}
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
	err := cmd.RunE(nil, []string{orchestratortest.MockCatalogID})
	assert.NoError(t, err)

	var response = &orchestrator.Catalog{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, orchestratortest.MockCatalogID, response.Id)
}

func TestNewGetCategoryCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetCategoryCommand()
	err := cmd.RunE(nil, []string{orchestratortest.MockCatalogID, orchestratortest.MockCategoryName})
	assert.NoError(t, err)

	var response = &orchestrator.Category{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, orchestratortest.MockCategoryName, response.Name)
}

func TestNewGetControlCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetControlCommand()
	err := cmd.RunE(nil, []string{orchestratortest.MockCatalogID, orchestratortest.MockCategoryName, orchestratortest.MockControlID})
	assert.NoError(t, err)

	var response = &orchestrator.Control{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, orchestratortest.MockControlID, response.Id)
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
