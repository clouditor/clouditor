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

	"clouditor.io/clouditor/api/orchestrator"
	apiruntime "clouditor.io/clouditor/api/runtime"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/persistence/inmemory"
	"clouditor.io/clouditor/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
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
		opts []ServiceOption
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "New service with database",
			args: args{
				opts: []ServiceOption{WithStorage(myStorage)},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, myStorage, s.storage)
			},
		},
		{
			name: "New service with catalogs file",
			args: args{
				opts: []ServiceOption{WithCatalogsFolder("catalogsFolder.json")},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, "catalogsFolder.json", s.catalogsFolder)
			},
		},
		{
			name: "New service with metrics file",
			args: args{
				opts: []ServiceOption{WithMetricsFile("metricsfile.json")},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, "metricsfile.json", s.metricsFile)
			},
		},
		{
			name: "New service with authorization strategy",
			args: args{
				opts: []ServiceOption{WithAuthorizationStrategy(&service.AuthorizationStrategyAllowAll{})},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, &service.AuthorizationStrategyAllowAll{}, s.authz)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			if tt.want != nil {
				tt.want(t, got, tt.args.opts)
			}
		})
	}
}

func TestCloudServiceHooks(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
		hookCounts      = 2
	)

	wg.Add(hookCounts)

	firstHookFunction := func(_ context.Context, cloudService *orchestrator.CloudService, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")
		wg.Done()
	}

	secondHookFunction := func(_ context.Context, cloudService *orchestrator.CloudService, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")
		wg.Done()
	}

	type args struct {
		in0               context.Context
		serviceUpdate     *orchestrator.UpdateCloudServiceRequest
		cloudServiceHooks []orchestrator.CloudServiceHookFunc
	}
	tests := []struct {
		name     string
		args     args
		wantResp *orchestrator.CloudService
		wantErr  bool
	}{
		{
			name: "Update Cloud Service",
			args: args{
				in0: context.TODO(),
				serviceUpdate: &orchestrator.UpdateCloudServiceRequest{
					CloudService: &orchestrator.CloudService{
						Id:          "00000000-0000-0000-0000-000000000000",
						Name:        "test service",
						Description: "test service",
					},
				},
				cloudServiceHooks: []orchestrator.CloudServiceHookFunc{firstHookFunction, secondHookFunction},
			},
			wantErr: false,
			wantResp: &orchestrator.CloudService{
				Id:          "00000000-0000-0000-0000-000000000000",
				Name:        "test service",
				Description: "test service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := NewService()

			_, err := s.CreateDefaultTargetCloudService()
			if err != nil {
				t.Errorf("CreateCloudService() error = %v", err)
			}

			for i, hookFunction := range tt.args.cloudServiceHooks {
				s.RegisterCloudServiceHook(hookFunction)

				// Check if hook is registered
				funcName1 := runtime.FuncForPC(reflect.ValueOf(s.cloudServiceHooks[i]).Pointer()).Name()
				funcName2 := runtime.FuncForPC(reflect.ValueOf(hookFunction).Pointer()).Name()
				assert.Equal(t, funcName1, funcName2)
			}

			// To test the hooks we have to call a function that calls the hook function
			gotResp, err := s.UpdateCloudService(tt.args.in0, tt.args.serviceUpdate)

			// wait for all hooks (2 services * 2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCloudService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotResp, tt.wantResp) {
				t.Errorf("UpdateCloudService() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			assert.Equal(t, tt.wantResp, gotResp)
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
		want    assert.ValueAssertionFunc
		wantErr bool
	}{
		{
			name: "return runtime",
			want: assert.NotNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			gotRes, err := svc.GetRuntimeInfo(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRuntimeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			tt.want(t, gotRes)
		})
	}
}

func TestService_ListPublicCertificates(t *testing.T) {
	type fields struct {
		UnimplementedOrchestratorServer orchestrator.UnimplementedOrchestratorServer
		cloudServiceHooks               []orchestrator.CloudServiceHookFunc
		toeHooks                        []orchestrator.TargetOfEvaluationHookFunc
		AssessmentResultHooks           []func(result *assessment.AssessmentResult, err error)
		storage                         persistence.Storage
		metricsFile                     string
		loadMetricsFunc                 func() ([]*assessment.Metric, error)
		catalogsFolder                  string
		loadCatalogsFunc                func() ([]*orchestrator.Catalog, error)
		events                          chan *orchestrator.MetricChangeEvent
		authz                           service.AuthorizationStrategy
	}
	type args struct {
		in0 context.Context
		req *orchestrator.ListPublicCertificatesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *orchestrator.ListPublicCertificatesResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Validation error",
			fields: fields{},
			args: args{
				req: nil,
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Pagination error",
			fields: fields{
				storage: &testutil.StorageWithError{ListErr: ErrSomeError},
			},
			args: args{
				req: &orchestrator.ListPublicCertificatesRequest{},
			},
			wantRes: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.Internal, status.Code(err))
				return assert.ErrorContains(t, err, "database error")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				storage: testutil.NewInMemoryStorage(t, func(s persistence.Storage) {
					// Create Certificate
					assert.NoError(t, s.Create(orchestratortest.NewCertificate()))
				}),
			},
			args: args{
				req: &orchestrator.ListPublicCertificatesRequest{},
			},
			wantRes: &orchestrator.ListPublicCertificatesResponse{
				Certificates: []*orchestrator.Certificate{
					{
						Id:             testdata.MockCertificateID,
						Name:           testdata.MockCertificateName,
						CloudServiceId: testdata.MockCloudServiceID1,
						IssueDate:      "2021-11-06",
						ExpirationDate: "2024-11-06",
						Standard:       testdata.MockCertificateName,
						AssuranceLevel: testdata.AssuranceLevelHigh,
						Cab:            testdata.MockCertificateCab,
						Description:    testdata.MockCertificateDescription,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				UnimplementedOrchestratorServer: tt.fields.UnimplementedOrchestratorServer,
				cloudServiceHooks:               tt.fields.cloudServiceHooks,
				toeHooks:                        tt.fields.toeHooks,
				AssessmentResultHooks:           tt.fields.AssessmentResultHooks,
				storage:                         tt.fields.storage,
				metricsFile:                     tt.fields.metricsFile,
				loadMetricsFunc:                 tt.fields.loadMetricsFunc,
				catalogsFolder:                  tt.fields.catalogsFolder,
				loadCatalogsFunc:                tt.fields.loadCatalogsFunc,
				events:                          tt.fields.events,
				authz:                           tt.fields.authz,
			}
			gotRes, err := svc.ListPublicCertificates(tt.args.in0, tt.args.req)
			assert.NoError(t, gotRes.Validate())

			tt.wantErr(t, err)

			if tt.wantRes != nil {
				if !reflect.DeepEqual(gotRes, tt.wantRes) {
					t.Errorf("Service.ListPublicCertificates() = %v, want %v", gotRes, tt.wantRes)
				}
			}
		})
	}
}
