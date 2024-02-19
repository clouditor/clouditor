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

package discovery

import (
	"bytes"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest"
	"clouditor.io/clouditor/v2/server"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestMain(m *testing.M) {
	svc := service_discovery.NewService()
	svc.StartDiscovery(&discoverytest.TestDiscoverer{TestCase: 2})

	os.Exit(clitest.RunCLITest(m,
		server.WithDiscovery(svc),
		server.WithExperimentalDiscovery(svc),
	))
}

func TestAddCommands(t *testing.T) {
	cmd := NewDiscoveryCommand()

	// Check if sub commands were added
	assert.True(t, cmd.HasSubCommands())

	// Check if NewQueryDiscoveryCommand was added
	for _, v := range cmd.Commands() {
		if v.Use == "query" {
			return
		}
	}
	t.Errorf("No query command was added")
}

func TestNewQueryDiscoveryCommandNoArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewQueryDiscoveryCommand()
	err = cmd.RunE(nil, []string{})
	assert.NoError(t, err)

	var response = &discovery.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}

func TestNewQueryDiscoveryCommandWithArgs(t *testing.T) {
	var err error
	var b bytes.Buffer

	cli.Output = &b

	cmd := NewQueryDiscoveryCommand()
	err = cmd.RunE(nil, []string{"Test Command"})
	assert.NoError(t, err)

	var response = &discovery.ListResourcesResponse{}
	err = protojson.Unmarshal(b.Bytes(), response)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}
