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

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/service"
	service_auth "clouditor.io/clouditor/service/auth"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	sock   net.Listener
	server *grpc.Server
)

func TestMain(m *testing.M) {
	var (
		err error
	)
	err = os.Chdir("../")
	if err != nil {
		panic(err)
	}

	s := service_orchestrator.NewService()
	sock, server, _, err = service.StartDedicatedAuthServer(":0", service_auth.WithApiKeySaveOnCreate(false))
	orchestrator.RegisterOrchestratorServer(server, s)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestSession(t *testing.T) {
	var (
		err     error
		session *Session
		dir     string
	)
	defer sock.Close()
	defer server.Stop()

	dir, err = ioutil.TempDir(os.TempDir(), ".clouditor")
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	viper.Set("session-directory", dir)

	session, err = NewSession(fmt.Sprintf("localhost:%d", sock.Addr().(*net.TCPAddr).Port))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer session.Close()

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, dir, session.Folder)

	client := auth.NewAuthenticationClient(session)

	var response *auth.TokenResponse

	// login with real user
	response, err = client.Login(context.Background(), &auth.LoginRequest{Username: "clouditor", Password: "clouditor"})

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)

	// update the session
	session.authorizer = api.NewInternalAuthorizerFromToken(
		session.authorizer.AuthURL(),
		&oauth2.Token{
			AccessToken:  response.AccessToken,
			TokenType:    response.TokenType,
			RefreshToken: response.RefreshToken,
			Expiry:       response.Expiry.AsTime(),
		},
	)

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
				if err != nil {
					return true
				} else {
					t.Errorf("Expected error.")
					return false
				}
			},
		},
		{
			name: "non-grpc error",
			args: args{
				msg: nil,
				err: fmt.Errorf("random error"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				if err != nil {
					return true
				} else {
					t.Errorf("Expected error.")
					return false
				}
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

// Test will fail due to no user input
func TestPromptForLogin(t *testing.T) {
	_, err := PromptForLogin()
	assert.Error(t, err)
}
