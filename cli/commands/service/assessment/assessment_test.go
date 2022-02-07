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

package assessment

import (
	"bytes"
	"context"
	"os"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/service"
	service_assessment "clouditor.io/clouditor/service/assessment"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	svc := service_assessment.NewService()

	r, _ := voc.ToStruct(&voc.Resource{ID: "myID", Type: []string{"Resource"}})

	// Add one evidence
	_, _ = svc.AssessEvidence(context.Background(), &assessment.AssessEvidenceRequest{
		Evidence: &evidence.Evidence{
			Id:        "00000000-0000-0000-000000000000",
			Resource:  r,
			ToolId:    "clouditor",
			Timestamp: timestamppb.Now(),
		},
	})

	os.Exit(clitest.RunCLITest(m, service.WithAssessment(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewAssessmentCommand()

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

func TestNewListStatisticsCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListStatisticsCommand()
	err := cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &assessment.ListStatisticsResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(1), response.NumberProcessedEvidences)
}
