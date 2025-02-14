// Copyright 2021-2022 Fraunhofer AISEC
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

package orchestrator

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	apiruntime "clouditor.io/clouditor/v2/api/runtime"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence/inmemory"
	"clouditor.io/clouditor/v2/service"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

func toStruct(f float32) (s *structpb.Value) {
	var (
		b   []byte
		err error
	)

	s = new(structpb.Value)

	b, err = json.Marshal(f)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(b, &s); err != nil {
		return nil
	}

	return
}

func TestNewService(t *testing.T) {
	var myStorage, err = inmemory.NewStorage()
	assert.NoError(t, err)

	type args struct {
		opts []service.Option[*Service]
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*Service]
	}{
		{
			name: "New service with database",
			args: args{
				opts: []service.Option[*Service]{WithStorage(myStorage)},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Same(t, myStorage, got.storage)
			},
		},
		{
			name: "New service with catalogs file",
			args: args{
				opts: []service.Option[*Service]{WithCatalogsFolder("catalogsFolder.json")},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "catalogsFolder.json", got.catalogsFolder)
			},
		},
		{
			name: "New service with metrics file",
			args: args{
				opts: []service.Option[*Service]{WithMetricsFile("metricsfile.json")},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "metricsfile.json", got.metricsFile)
			},
		},
		{
			name: "New service with authorization strategy",
			args: args{
				opts: []service.Option[*Service]{WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{})},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal[service.AuthorizationStrategy](t, &service.AuthorizationStrategyAllowAll{}, got.authz)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)
			tt.want(t, got)
		})
	}
}

func TestCertificationTargetHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
		hookCounts      = 2
	)

	wg.Add(hookCounts)

	firstHookFunction := func(_ context.Context, CertificationTarget *orchestrator.CertificationTarget, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")
		wg.Done()
	}

	secondHookFunction := func(_ context.Context, CertificationTarget *orchestrator.CertificationTarget, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")
		wg.Done()
	}

	type args struct {
		in0                      context.Context
		targetUpdate             *orchestrator.UpdateCertificationTargetRequest
		CertificationTargetHooks []orchestrator.CertificationTargetHookFunc
	}
	tests := []struct {
		name    string
		args    args
		wantRes assert.Want[*orchestrator.CertificationTarget]
		wantErr assert.WantErr
	}{
		{
			name: "Update Certification Target",
			args: args{
				in0: context.TODO(),
				targetUpdate: &orchestrator.UpdateCertificationTargetRequest{
					CertificationTarget: &orchestrator.CertificationTarget{
						Id:          "00000000-0000-0000-0000-000000000000",
						Name:        "test target",
						Description: "test target",
					},
				},
				CertificationTargetHooks: []orchestrator.CertificationTargetHookFunc{firstHookFunction, secondHookFunction},
			},
			wantErr: assert.Nil[error],
			wantRes: func(t *testing.T, got *orchestrator.CertificationTarget) bool {
				return assert.Equal(t, "00000000-0000-0000-0000-000000000000", got.Id) &&
					assert.Equal(t, "test target", got.Name) && assert.Equal(t, "test target", got.Description)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService()

			_, err := s.CreateDefaultCertificationTarget()
			if err != nil {
				t.Errorf("CreateCertificationTarget() error = %v", err)
			}

			for i, hookFunction := range tt.args.CertificationTargetHooks {
				s.CreateCertificationTargetHook(hookFunction)

				// Check if hook is registered
				funcName1 := runtime.FuncForPC(reflect.ValueOf(s.CertificationTargetHooks[i]).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(hookFunction).Pointer()).Name()
				assert.Equal(t, funcName1, funcName2)
			}

			// To test the hooks we have to call a function that calls the hook function
			gotRes, err := s.UpdateCertificationTarget(tt.args.in0, tt.args.targetUpdate)

			// wait for all hooks (2 certification targets * 2 hooks)
			wg.Wait()

			tt.wantErr(t, err)
			tt.wantRes(t, gotRes)

			assert.Equal(t, hookCounts, hookCallCounter)
		})
	}
}

func TestService_GetRuntimeInfo(t *testing.T) {
	type fields struct {
	}
	type args struct {
		in0 context.Context
		in1 *apiruntime.GetRuntimeInfoRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*apiruntime.Runtime]
		wantErr assert.WantErr
	}{
		{
			name:    "return runtime",
			want:    assert.NotNil[*apiruntime.Runtime],
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			gotRes, err := svc.GetRuntimeInfo(tt.args.in0, tt.args.in1)

			tt.wantErr(t, err)
			tt.want(t, gotRes)
		})
	}
}

func TestDefaultServiceSpec(t *testing.T) {
	tests := []struct {
		name string
		want assert.Want[launcher.ServiceSpec]
	}{
		{
			name: "Happy path",
			want: func(t *testing.T, got launcher.ServiceSpec) bool {
				return assert.NotNil(t, got)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultServiceSpec()

			tt.want(t, got)
		})
	}
}
