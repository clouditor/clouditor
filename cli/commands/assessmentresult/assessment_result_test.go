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

package assessmentresult

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
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
		Result: clitest.MockAssessmentResult1})
	if err != nil {
		panic(err)
	}

	os.Exit(clitest.RunCLITest(m, server.WithServices(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewAssessmentResultCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewListAssessmentResultsCommand was added
	for _, v := range cmd.Commands() {
		if v.Use == "list" {
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
