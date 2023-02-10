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
	"strings"
	"sync"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence/inmemory"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, myStorage, service.storage)
			},
		},
		{
			name: "New service with catalogs file",
			args: args{
				opts: []ServiceOption{WithCatalogsFile("catalogsfile.json")},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, "catalogsfile.json", service.catalogsFile)
			},
		},
		{
			name: "New service with metrics file",
			args: args{
				opts: []ServiceOption{WithMetricsFile("metricsfile.json")},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				service, ok := i1.(*Service)
				if !assert.True(tt, ok) {
					return false
				}

				return assert.Equal(tt, "metricsfile.json", service.metricsFile)
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

func Test_CreateCertificate(t *testing.T) {
	// Mock certificates
	mockCertificate := orchestratortest.NewCertificate()
	mockCertificateWithoutID := orchestratortest.NewCertificate()
	mockCertificateWithoutID.Id = ""

	type args struct {
		in0 context.Context
		req *orchestrator.CreateCertificateRequest
	}
	var tests = []struct {
		name         string
		args         args
		wantResponse *orchestrator.Certificate
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			"missing request",
			args{
				context.Background(),
				nil,
			},
			nil,
			func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			"missing certificate",
			args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{},
			},
			nil,
			func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "Certificate: value is required")
			},
		},
		{
			"missing certificate id",
			args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificateWithoutID,
				},
			},
			nil,
			func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "Id: value length must be at least 1 runes")
			},
		},
		{
			"valid certificate",
			args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: mockCertificate,
				},
			},
			mockCertificate,
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResponse, err := s.CreateCertificate(tt.args.in0, tt.args.req)
			assert.NoError(t, gotResponse.Validate())

			tt.wantErr(t, err)

			// If no error is wanted, check response
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Service.CreateCertificate() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func Test_UpdateCertificate(t *testing.T) {
	var (
		certificate *orchestrator.Certificate
		err         error
	)
	orchestratorService := NewService()

	// 1st case: Certificate is nil
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 2nd case: Certificate ID is nil
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: certificate,
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	// 3rd case: Certificate not found since there are no certificates yet
	_, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: &orchestrator.Certificate{
			Id:             testdata.MockCertificateID,
			Name:           "EUCS",
			CloudServiceId: testdata.MockCloudServiceID,
		},
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Certificate updated successfully
	mockCertificate := orchestratortest.NewCertificate()
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	// update the certificate's description and send the update request
	mockCertificate.Description = "new description"
	certificate, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		Certificate: mockCertificate,
	})
	assert.NoError(t, certificate.Validate())
	assert.NoError(t, err)
	assert.NotNil(t, certificate)
	assert.Equal(t, "new description", certificate.Description)
}

func Test_RemoveCertificate(t *testing.T) {
	var (
		err                      error
		listCertificatesResponse *orchestrator.ListCertificatesResponse
	)
	orchestratorService := NewService()

	// 1st case: Empty certificate ID error
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: ""})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.InvalidArgument)

	// 2nd case: ErrRecordNotFound
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: "0000"})
	assert.Error(t, err)
	assert.Equal(t, status.Code(err), codes.NotFound)

	// 3rd case: Record removed successfully
	mockCertificate := orchestratortest.NewCertificate()
	assert.NoError(t, mockCertificate.Validate())
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)

	// There is a record for certificates in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.NotEmpty(t, listCertificatesResponse.Certificates)

	// Remove record
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: mockCertificate.Id})
	assert.NoError(t, err)

	// There is a record for cloud services in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.Empty(t, listCertificatesResponse.Certificates)
}

func Test_GetCertificate(t *testing.T) {
	tests := []struct {
		name    string
		req     *orchestrator.GetCertificateRequest
		res     *orchestrator.Certificate
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "invalid request",
			req:  nil,
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "certificate not found",
			req:  &orchestrator.GetCertificateRequest{CertificateId: ""},
			res:  nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
				return assert.ErrorContains(t, err, "invalid request: invalid GetCertificateRequest.CertificateId: value length must be at least 1 runes")
			},
		},
		{
			name:    "valid",
			req:     &orchestrator.GetCertificateRequest{CertificateId: testdata.MockCertificateID},
			res:     orchestratortest.NewCertificate(),
			wantErr: assert.NoError,
		},
	}
	orchestratorService := NewService()

	// Create Certificate
	if err := orchestratorService.storage.Create(orchestratortest.NewCertificate()); err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.GetCertificate(context.Background(), tt.req)
			assert.NoError(t, res.Validate())

			tt.wantErr(t, err)

			if tt.res != nil {
				assert.NotEmpty(t, res.Id)
				// Compare timestamp. We have to cut off the microseconds, otherwise an error is returned.
				tt.res.States[0].Timestamp = strings.Split(tt.res.States[0].GetTimestamp(), ".")[0]
				res.States[0].Timestamp = strings.Split(res.States[0].GetTimestamp(), ".")[0]
				assert.True(t, proto.Equal(tt.res, res), "Want: %v\nGot : %v", tt.res, res)
			}
		})
	}
}

func Test_ListCertificates(t *testing.T) {
	var (
		listCertificatesResponse *orchestrator.ListCertificatesResponse
		err                      error
	)

	orchestratorService := NewService()
	// 1st case: No services stored
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.Empty(t, listCertificatesResponse.Certificates)

	// 2nd case: One service stored
	err = orchestratorService.storage.Create(orchestratortest.NewCertificate())
	assert.NoError(t, err)

	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	// We check only the first certificate and assume that all certificates are valid
	assert.NoError(t, listCertificatesResponse.Certificates[0].Validate())
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.NotEmpty(t, listCertificatesResponse.Certificates)
	assert.Equal(t, len(listCertificatesResponse.Certificates), 1)

	// 3rd case: Invalid request
	_, err = orchestratorService.ListCertificates(context.Background(),
		&orchestrator.ListCertificatesRequest{OrderBy: "not a field"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, err.Error(), api.ErrInvalidColumnName.Error())
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
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("UpdateCloudService() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			assert.Equal(t, tt.wantResp, gotResp)
			assert.Equal(t, hookCounts, hookCallCounter)
		})
	}
}

/*
func TestService_Runtime(t *testing.T) {
	type fields struct {
		storage persistence.Storage
		authz   service.AuthorizationStrategy
		runtime *orchestrator.Runtime
	}
	type args struct {
		in0 context.Context
		in1 *orchestrator.RuntimeRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRuntime *orchestrator.RuntimeResponse
		wantErr     bool
	}{
		{
			name: "Happy path",
			fields: fields{
				runtime: &orchestrator.Runtime{
					ReleaseVersion: "testVersion",
				},
			},
			wantRuntime: &orchestrator.RuntimeResponse{
				Runtime: &orchestrator.Runtime{
					ReleaseVersion: "testVersion",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				storage: tt.fields.storage,
				authz:   tt.fields.authz,
			}
			gotRuntime, err := svc.Runtime(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Runtime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRuntime, tt.wantRuntime) {
				t.Errorf("Service.Runtime() = %v, want %v", gotRuntime, tt.wantRuntime)
			}
		})
	}
}
*/
