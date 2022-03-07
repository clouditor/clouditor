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

package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/service"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	sock    net.Listener
	srv     *grpc.Server
	authSrv *oauth2.AuthorizationServer
)

func TestMain(m *testing.M) {
	var (
		err  error
		port int
		code int
	)

	// Because of a cyclic dependency, we cannot use clitest here
	err = os.Chdir("../")
	if err != nil {
		panic(err)
	}

	authSrv, port, err = testutil.StartAuthenticationServer()
	if err != nil {
		panic(err)
	}

	s := service_orchestrator.NewService()
	sock, srv, err = service.StartGRPCServer(testutil.JWKSURL(port), service.WithOrchestrator(s))
	if err != nil {
		panic(err)
	}

	code = m.Run()

	authSrv.Close()
	srv.Stop()

	os.Exit(code)
}

func TestSession(t *testing.T) {
	var (
		err     error
		session *Session
		dir     string
		token   *oauth2.Token
	)

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	viper.Set("session-directory", dir)

	token, err = authSrv.GenerateToken(testutil.TestAuthClientID, 0, 0)

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.AccessToken)

	session, err = NewSession(fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port), &oauth2.Config{
		ClientID: testutil.TestAuthClientID,
		Endpoint: oauth2.Endpoint{},
	}, token)

	assert.NoError(t, err)

	err = session.Save()

	assert.NoError(t, err)

	session, err = ContinueSession()
	assert.NoError(t, err)
	assert.NotNil(t, session)

	// Do a simple authenticated call
	oc := orchestrator.NewOrchestratorClient(session)
	_, err = oc.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{})
	assert.NoError(t, err)
}

func TestSession_HandleResponse(t *testing.T) {
	type fields struct {
		URL        string
		Token      oauth2.Token
		Folder     string
		ClientConn *grpc.ClientConn
	}
	type args struct {
		msg proto.Message
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "grpc Error",
			args: args{
				msg: nil,
				err: status.Errorf(codes.Internal, "internal error occurred!"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
		{
			name: "non-grpc error",
			args: args{
				msg: nil,
				err: fmt.Errorf("random error"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				URL: tt.fields.URL,
				//Token:      tt.fields.Token,
				Folder:     tt.fields.Folder,
				ClientConn: tt.fields.ClientConn,
			}
			tt.wantErr(t, s.HandleResponse(tt.args.msg, tt.args.err), fmt.Sprintf("HandleResponse(%v, %v)", tt.args.msg, tt.args.err))
		})
	}
}
