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

package server_test

import (
	"context"
	"os"
	"testing"

	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
)

func TestMain(m *testing.M) {
	svc := service_orchestrator.NewService()

	os.Exit(clitest.RunCLITest(m, server.WithServices(svc), server.WithReflection()))
}

func TestReflectionNoAuth(t *testing.T) {
	var (
		session *cli.Session
		conn    *grpc.ClientConn
		client  grpc_reflection_v1.ServerReflectionClient
		sclient grpc_reflection_v1.ServerReflection_ServerReflectionInfoClient
		res     *grpc_reflection_v1.ServerReflectionResponse
		err     error
	)

	session, err = cli.ContinueSession()
	assert.NoError(t, err)

	// Only use the host from the session, but not the (authentication) connection, since we want to test, whether we
	// can access the reflection without authentication
	conn, err = grpc.NewClient(session.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	client = grpc_reflection_v1.NewServerReflectionClient(conn)
	sclient, err = client.ServerReflectionInfo(context.TODO())
	assert.NoError(t, err)

	err = sclient.Send(&grpc_reflection_v1.ServerReflectionRequest{
		Host:           "localhost",
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_ListServices{},
	})
	assert.NoError(t, err)

	res, err = sclient.Recv()
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
