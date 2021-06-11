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

package assessment_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/cli"
	cli_assessment "clouditor.io/clouditor/cli/commands/assessment"
	"clouditor.io/clouditor/cli/commands/login"
	service_assessment "clouditor.io/clouditor/service/assessment"
	service_auth "clouditor.io/clouditor/service/auth"
	"clouditor.io/clouditor/service/standalone"
	"clouditor.io/clouditor/voc"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

var sock net.Listener
var server *grpc.Server

func TestMain(m *testing.M) {
	var (
		dir string
		err error
	)

	// make sure, that we are in the clouditor root folder to find the policies
	err = os.Chdir("../../..")
	if err != nil {
		panic(err)
	}

	var ready chan bool = make(chan bool)

	sock, server, err = service_auth.StartDedicatedAuthServer(":0")
	if err != nil {
		panic(err)
	}

	defer sock.Close()
	defer server.Stop()

	assessmentServer := standalone.NewAssessmentServer().(*service_assessment.Service)
	assessmentServer.ResultHook = func(result *assessment.Result, err error) {
		ready <- true
	}
	assessment.RegisterAssessmentServer(server, assessmentServer)

	client := standalone.NewAssessmentClient()

	resource := &voc.ObjectStorageResource{
		StorageResource: voc.StorageResource{
			Resource: voc.Resource{
				ID:   "some-id",
				Type: []string{"ObjectStorage", "Storage", "Resource"},
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			TransportEncryption: voc.NewTransportEncryption(true, true, "TLS1_2"),
		},
	}

	s, err := voc.ToStruct(resource)
	if err != nil {
		panic(err)
	}

	evidence := &assessment.Evidence{
		ResourceId:        "some-id",
		ApplicableMetrics: []int32{1},
		Resource:          s,
	}

	_, err = client.StoreEvidence(context.Background(), &assessment.StoreEvidenceRequest{
		Evidence: evidence,
	})
	if err != nil {
		panic(err)
	}

	// make the test wait for envidence to be stored
	select {
	case <-ready:
		break
	case <-time.After(10 * time.Second):
		panic("Timeout while waiting for evidence assessment result to be ready")
	}

	// TODO(oxisto): refactor the next lines into a common test login function
	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	if err != nil {
		panic(err)
	}

	viper.Set("username", "clouditor")
	viper.Set("password", "clouditor")
	viper.Set("session-directory", dir)

	cmd := login.NewLoginCommand()
	err = cmd.RunE(nil, []string{fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port)})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestListResults(t *testing.T) {
	var b bytes.Buffer
	var err error

	cli.Output = &b

	cmd := cli_assessment.NewListResultsCommand()
	err = cmd.RunE(nil, []string{})

	assert.Nil(t, err)

	var response *assessment.ListAssessmentResultsResponse = &assessment.ListAssessmentResultsResponse{}

	err = protojson.Unmarshal(b.Bytes(), response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Results)
}
