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

package evidence

import (
	"bytes"
	"context"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_evidence "clouditor.io/clouditor/v2/service/evidence"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	var (
		svc *service_evidence.Service
		err error
	)

	svc = service_evidence.NewService()

	_, err = svc.StoreEvidence(context.TODO(), &evidence.StoreEvidenceRequest{Evidence: &evidence.Evidence{
		Id:                   testdata.MockTargetOfEvaluationID1,
		TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
		ToolId:               testdata.MockEvidenceToolID1,
		Timestamp:            timestamppb.Now(),
		Resource:             ontology.ProtoResource(&ontology.VirtualMachine{Id: testdata.MockVirtualMachineID1, Name: testdata.MockVirtualMachineName1}),
	}})
	if err != nil {
		panic(err)
	}

	os.Exit(clitest.RunCLITest(m, server.WithServices(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewEvidenceCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewListResultsCommand was added
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

	cmd := NewListEvidencesCommand()
	err := cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &evidence.ListEvidencesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Evidences)
}
