// Copyright 2022 Fraunhofer AISEC
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

package service_test

import (
	"context"
	"os"
	"testing"

	"clouditor.io/clouditor/cli"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/service"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func TestMain(m *testing.M) {
	svc := service_orchestrator.NewService()

	os.Exit(clitest.RunCLITest(m, service.WithOrchestrator(svc)))
}

func TestReflectionNoAuth(t *testing.T) {
	var (
		session *cli.Session
		conn    *grpc.ClientConn
		client  grpc_reflection_v1alpha.ServerReflectionClient
		sclient grpc_reflection_v1alpha.ServerReflection_ServerReflectionInfoClient
		res     *grpc_reflection_v1alpha.ServerReflectionResponse
		err     error
	)

	session, err = cli.ContinueSession()
	assert.NoError(t, err)

	// Only use the host from the session, but not the (authentication) connection, since we want to test, whether we
	// can access the reflection without authentication
	conn, err = grpc.Dial(session.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	client = grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	sclient, err = client.ServerReflectionInfo(context.TODO())
	assert.NoError(t, err)

	err = sclient.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host:           "localhost",
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{},
	})
	assert.NoError(t, err)

	res, err = sclient.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
