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
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/orchestratortest"
	"clouditor.io/clouditor/persistence/inmemory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	os.Exit(m.Run())
}

func TestAssessmentResultHook(t *testing.T) {
	var (
		hookCallCounter = 0
		wg              sync.WaitGroup
	)
	wg.Add(2)

	firstHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		wg.Done()
	}

	secondHookFunction := func(assessmentResult *assessment.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")

		wg.Done()
	}

	service := NewService()
	service.RegisterAssessmentResultHook(firstHookFunction)
	service.RegisterAssessmentResultHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultHooks[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultHooks[1]).Pointer()).Name()
	funcName2 = runtime.FuncForPC(reflect.ValueOf(secondHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check GRPC call
	type args struct {
		in0        context.Context
		assessment *orchestrator.StoreAssessmentResultRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *orchestrator.StoreAssessmentResultResponse
		wantErr  bool
	}{
		{
			name: "Store first assessment result to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						MetricId:       testdata.MockMetricID,
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       ">=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)

			// wait for all hooks (2 hooks)
			wg.Wait()

			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResult() error = %v, wantErrMessage %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreAssessmentResult() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.NotEmpty(t, s.results)
			assert.Equal(t, 2, hookCallCounter)
		})
	}
}

func TestStoreAssessmentResult(t *testing.T) {
	type args struct {
		in0        context.Context
		assessment *orchestrator.StoreAssessmentResultRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *orchestrator.StoreAssessmentResultResponse
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Store assessment to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						MetricId:       "assessmentResultMetricID",
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       "<=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr:  assert.NoError,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
		{
			name: "Store assessment without metricId to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:             testdata.MockAssessmentResultID,
						EvidenceId:     testdata.MockEvidenceID,
						CloudServiceId: testdata.MockCloudServiceID,
						Timestamp:      timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue:    toStruct(1.0),
							Operator:       "<=",
							IsDefault:      true,
							CloudServiceId: testdata.MockCloudServiceID,
							MetricId:       testdata.MockMetricID,
						},
						NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
						Compliant:             true,
						ResourceId:            testdata.MockResourceID,
						ResourceTypes:         []string{"ResourceType"},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "caused by: invalid AssessmentResult.MetricId: value length must be at least 1 runes")
			},
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantResp, gotResp)

			if err == nil {
				assert.NotNil(t, s.results[testdata.MockAssessmentResultID])
			} else {
				assert.Empty(t, s.results)
			}
		})
	}
}

func TestStoreAssessmentResults(t *testing.T) {
	const (
		count1 = 1
		count2 = 2
	)

	type fields struct {
		countElementsInMock    int
		countElementsInResults int
	}

	type args struct {
		streamToServer            *mockStreamer
		streamToClientWithSendErr *mockStreamerWithSendErr
		streamToServerWithRecvErr *mockStreamerWithRecvErr
	}

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantRespMessage []orchestrator.StoreAssessmentResultResponse
		wantErrMessage  string
	}{
		{
			name: "Store 2 assessment results to the map",
			fields: fields{
				countElementsInMock:    count2,
				countElementsInResults: count2,
			},
			args:    args{streamToServer: createMockStream(createStoreAssessmentResultRequestsMock(count2))},
			wantErr: false,
			wantRespMessage: []orchestrator.StoreAssessmentResultResponse{
				{
					Status: true,
				},
				{
					Status: true,
				},
			},
		},
		{
			name: "Missing MetricID",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:    args{streamToServer: createMockStream(createStoreAssessmentResultRequestMockWithMissingMetricID(count1))},
			wantErr: false,
			wantRespMessage: []orchestrator.StoreAssessmentResultResponse{
				{
					Status:        false,
					StatusMessage: "MetricId: value length must be at least 1 runes",
				},
			},
		},
		{
			name: "Error in stream to server - Recv()-err",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:           args{streamToServerWithRecvErr: createMockStreamWithRecvErr(createStoreAssessmentResultRequestsMock(count1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot receive stream request",
		},
		{
			name: "Error in stream to client - Send()-err",
			fields: fields{
				countElementsInMock:    count1,
				countElementsInResults: 0,
			},
			args:           args{streamToClientWithSendErr: createMockStreamWithSendErr(createStoreAssessmentResultRequestsMock(count1))},
			wantErr:        true,
			wantErrMessage: "rpc error: code = Unknown desc = cannot stream response to the client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()

			var err error

			if tt.args.streamToServer != nil {
				err = s.StoreAssessmentResults(tt.args.streamToServer)
			} else if tt.args.streamToClientWithSendErr != nil {
				err = s.StoreAssessmentResults(tt.args.streamToClientWithSendErr)
			} else if tt.args.streamToServerWithRecvErr != nil {
				err = s.StoreAssessmentResults(tt.args.streamToServerWithRecvErr)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Got StoreAssessmentResults() error = %v, wantErrMessage %v", err, tt.wantErr)
				assert.Equal(t, tt.fields.countElementsInResults, len(s.results))
				return
			} else if tt.wantErr {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			} else {
				// Close stream for testing
				close(tt.args.streamToServer.SentFromServer)
				assert.Nil(t, err)
				assert.Equal(t, tt.fields.countElementsInResults, len(s.results))

				// Check all stream responses from server to client
				i := 0
				for elem := range tt.args.streamToServer.SentFromServer {
					assert.Contains(t, elem.StatusMessage, tt.wantRespMessage[i].StatusMessage)
					assert.Equal(t, tt.wantRespMessage[i].Status, elem.Status)
					i++
				}
			}

		})
	}
}

// createStoreAssessmentResultRequestMockWithMissingMetricID create one StoreAssessmentResultRequest without ToolID
func createStoreAssessmentResultRequestMockWithMissingMetricID(count int) []*orchestrator.StoreAssessmentResultRequest {
	req := createStoreAssessmentResultRequestsMock(count)

	req[0].Result.MetricId = ""

	return req
}

// createStoreAssessmentResultrequestMocks creates store assessment result requests with random assessment result IDs
func createStoreAssessmentResultRequestsMock(count int) []*orchestrator.StoreAssessmentResultRequest {
	var mockRequests []*orchestrator.StoreAssessmentResultRequest

	for i := 0; i < count; i++ {
		storeAssessmentResultRequest := &orchestrator.StoreAssessmentResultRequest{
			Result: &assessment.AssessmentResult{
				Id:             uuid.NewString(),
				MetricId:       fmt.Sprintf("assessmentResultMetricID-%d", i),
				EvidenceId:     testdata.MockEvidenceID,
				CloudServiceId: testdata.MockCloudServiceID,
				Timestamp:      timestamppb.Now(),
				MetricConfiguration: &assessment.MetricConfiguration{
					TargetValue:    toStruct(1.0),
					Operator:       "<=",
					IsDefault:      true,
					CloudServiceId: testdata.MockCloudServiceID,
					MetricId:       testdata.MockMetricID,
				},
				NonComplianceComments: testdata.MockAssessmentResultNonComplianceComment,
				Compliant:             true,
				ResourceId:            testdata.MockResourceID,
				ResourceTypes:         []string{"ResourceType"},
			},
		}

		mockRequests = append(mockRequests, storeAssessmentResultRequest)
	}

	return mockRequests
}

type mockStreamer struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultResponse
}

func createMockStream(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamer {
	m := &mockStreamer{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultResponse, len(requests))
	return m
}

func (m mockStreamer) Send(response *orchestrator.StoreAssessmentResultResponse) error {
	m.SentFromServer <- response
	return nil
}

func (m mockStreamer) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func (mockStreamer) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (mockStreamer) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (mockStreamer) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (mockStreamer) Context() context.Context {
	panic("implement me")
}

func (mockStreamer) SendMsg(_ interface{}) error {
	panic("implement me")
}

func (mockStreamer) RecvMsg(_ interface{}) error {
	panic("implement me")
}

type mockStreamerWithSendErr struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultResponse
}

func (mockStreamerWithSendErr) Send(*orchestrator.StoreAssessmentResultResponse) error {
	return errors.New("Send()-err")
}

func (m mockStreamerWithSendErr) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	if len(m.RecvToServer) == 0 {
		return nil, io.EOF
	}
	req, more := <-m.RecvToServer
	if !more {
		return nil, errors.New("empty")
	}

	return req, nil
}

func createMockStreamWithSendErr(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamerWithSendErr {
	m := &mockStreamerWithSendErr{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultResponse, len(requests))
	return m
}

type mockStreamerWithRecvErr struct {
	grpc.ServerStream
	RecvToServer   chan *orchestrator.StoreAssessmentResultRequest
	SentFromServer chan *orchestrator.StoreAssessmentResultResponse
}

func (mockStreamerWithRecvErr) Send(*orchestrator.StoreAssessmentResultResponse) error {
	panic("implement me")
}

func (mockStreamerWithRecvErr) Recv() (*orchestrator.StoreAssessmentResultRequest, error) {
	err := errors.New("Recv()-error")

	return nil, err
}

func createMockStreamWithRecvErr(requests []*orchestrator.StoreAssessmentResultRequest) *mockStreamerWithRecvErr {
	m := &mockStreamerWithRecvErr{
		RecvToServer: make(chan *orchestrator.StoreAssessmentResultRequest, len(requests)),
	}
	for _, req := range requests {
		m.RecvToServer <- req
	}

	m.SentFromServer = make(chan *orchestrator.StoreAssessmentResultResponse, len(requests))
	return m
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
