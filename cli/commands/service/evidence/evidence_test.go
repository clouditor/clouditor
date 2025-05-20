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
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_evidence "clouditor.io/clouditor/v2/service/evidence"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestMain(m *testing.M) {
	svc := service_evidence.NewService()
	svc.StoreEvidence(context.Background(), &evidence.StoreEvidenceRequest{
		Evidence: clitest.MockEvidence1,
	})
	// svc.StartDiscovery(&discoverytest.TestDiscoverer{TestCase: 2})

	os.Exit(clitest.RunCLITest(m,
		server.WithServices(svc),
	))
}

func TestAddCommands(t *testing.T) {
	cmd := NewEvidenceCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())
}

func TestNewListResourcesCommandNoArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListResourceCommand()
	err = cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &evidence.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}

func TestNewListResourcesCommandWithArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListResourceCommand()
	err = cmd.RunE(nil, []string{"Test Command"})
	assert.NoError(t, err)

	var response = &evidence.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNewListEvidencesCommandNoArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListEvidencesCommand()
	err = cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &evidence.ListEvidencesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNewListEvidencesCommandWithArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListEvidencesCommand()
	err = cmd.RunE(nil, []string{"Test Command"})
	assert.NoError(t, err)

	var response = &evidence.ListEvidencesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestNewGetEvidenceCommand(t *testing.T) {
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewGetEvidenceCommand()
	err := cmd.RunE(nil, []string{clitest.MockEvidence1.Id})
	assert.NoError(t, err)

	var response = &evidence.Evidence{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response)
	assert.Equal(t, clitest.MockEvidence1.Id, response.Id)
}
