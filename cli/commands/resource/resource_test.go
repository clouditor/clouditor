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

package resource

import (
	"bytes"
	"os"
	"testing"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/servicetest/discoverytest"
	"clouditor.io/clouditor/server"
	service_discovery "clouditor.io/clouditor/service/discovery"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestMain(m *testing.M) {
	svc := service_discovery.NewService()
	svc.StartDiscovery(&discoverytest.TestDiscoverer{TestCase: 2})

	os.Exit(clitest.RunCLITest(m, server.WithDiscovery(svc)))
}

func TestAddCommands(t *testing.T) {
	cmd := NewResourceCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewListCommand was added
	for _, v := range cmd.Commands() {
		if v.Use == "list" {
			return
		}
	}
	t.Errorf("No list command was added")
}

func TestNewListCommand(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewListResourcesCommand()
	err = cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &discovery.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}
