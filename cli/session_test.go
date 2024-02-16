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

package cli_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/orchestratortest"
	"clouditor.io/clouditor/v2/server"
	service_orchestrator "clouditor.io/clouditor/v2/service/orchestrator"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	var (
		svc *service_orchestrator.Service
		err error
	)

	clitest.AutoChdir()

	svc = service_orchestrator.NewService()
	_, err = svc.CreateDefaultTargetCloudService()
	if err != nil {
		panic(err)
	}

	_, err = svc.CreateCatalog(context.TODO(), &orchestrator.CreateCatalogRequest{Catalog: orchestratortest.NewCatalog()})
	if err != nil {
		panic(err)
	}

	os.Exit(clitest.RunCLITest(m, server.WithOrchestrator(svc)))
}

func TestSession(t *testing.T) {
	var (
		err     error
		session *cli.Session
	)

	session, err = cli.ContinueSession()
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
			s := &cli.Session{
				URL:        tt.fields.URL,
				Folder:     tt.fields.Folder,
				ClientConn: tt.fields.ClientConn,
			}
			tt.wantErr(t, s.HandleResponse(tt.args.msg, tt.args.err), fmt.Sprintf("HandleResponse(%v, %v)", tt.args.msg, tt.args.err))
		})
	}
}

func TestValidArgsGetMetrics(t *testing.T) {
	type args struct {
		in0        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  assert.ValueAssertionFunc
		want1 cobra.ShellCompDirective
	}{
		{
			name: "some metrics",
			args: args{
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.NotNil(tt, i1)
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cli.ValidArgsGetMetrics(tt.args.in0, tt.args.args, tt.args.toComplete)

			tt.want(t, got)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidArgsGetMetrics() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidArgsGetCatalogs(t *testing.T) {
	type args struct {
		in0        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  assert.ValueAssertionFunc
		want1 cobra.ShellCompDirective
	}{
		{
			name: "some catalogs",
			args: args{
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s: %s", testdata.MockCatalogID, testdata.MockCatalogName, testdata.MockCatalogDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "all args - return nothing",
			args: args{
				args:       []string{testdata.MockCatalogID},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Empty(t, i1)
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cli.ValidArgsGetCatalogs(tt.args.in0, tt.args.args, tt.args.toComplete)

			tt.want(t, got)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidArgsGetCatalogs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidArgsGetCategory(t *testing.T) {
	type args struct {
		in0        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  assert.ValueAssertionFunc
		want1 cobra.ShellCompDirective
	}{
		{
			name: "no arg - return catalog",
			args: args{
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s: %s", testdata.MockCatalogID, testdata.MockCatalogName, testdata.MockCatalogDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "one arg - return category",
			args: args{
				args:       []string{testdata.MockCatalogID},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s", testdata.MockCategoryName, testdata.MockCategoryDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "all args - return nothing",
			args: args{
				args:       []string{testdata.MockCatalogID, testdata.MockCategoryName},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Empty(t, i1)
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cli.ValidArgsGetCategory(tt.args.in0, tt.args.args, tt.args.toComplete)

			tt.want(t, got)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidArgsGetCatalogs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidArgsGetControls(t *testing.T) {
	type args struct {
		in0        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  assert.ValueAssertionFunc
		want1 cobra.ShellCompDirective
	}{
		{
			name: "no arg - return catalog",
			args: args{
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s: %s", testdata.MockCatalogID, testdata.MockCatalogName, testdata.MockCatalogDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "one arg - return category",
			args: args{
				args:       []string{testdata.MockCatalogID},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s", testdata.MockCategoryName, testdata.MockCategoryDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "two args - return category",
			args: args{
				args:       []string{testdata.MockCatalogID, testdata.MockCategoryName},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Contains(t, i1, fmt.Sprintf("%s\t%s: %s", testdata.MockControlID1, testdata.MockControlName, testdata.MockControlDescription))
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
		{
			name: "all args - return nothing",
			args: args{
				args:       []string{testdata.MockCatalogID, testdata.MockCategoryName, testdata.MockControlID1},
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Empty(t, i1)
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cli.ValidArgsGetControls(tt.args.in0, tt.args.args, tt.args.toComplete)

			tt.want(t, got)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ValidArgsGetCatalogs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidArgsGetCloudServices(t *testing.T) {
	type args struct {
		in0        *cobra.Command
		args       []string
		toComplete string
	}
	tests := []struct {
		name  string
		args  args
		want  assert.ValueAssertionFunc
		want1 cobra.ShellCompDirective
	}{
		{
			name: "some cloud services",
			args: args{
				toComplete: "",
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.NotNil(tt, i1)
			},
			want1: cobra.ShellCompDirectiveNoFileComp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cli.ValidArgsGetCloudServices(tt.args.in0, tt.args.args, tt.args.toComplete)

			tt.want(t, got)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TestValidArgsGetCloudServices() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
