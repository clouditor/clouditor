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
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/internal/testutil/clitest"
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

const (
	assessmentResultID1 = "11111111-1111-1111-1111-111111111111"
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
						Id:         assessmentResultID1,
						MetricId:   "assessmentResultMetricID",
						EvidenceId: "11111111-1111-1111-1111-111111111111",
						Timestamp:  timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue: toStruct(1.0),
							Operator:    "operator",
							IsDefault:   true,
						},
						NonComplianceComments: "non_compliance_comment",
						Compliant:             true,
						ResourceId:            "resourceID",
					},
				},
			},
			wantErr: false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{
				Status: true,
			},
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
				t.Errorf("StoreAssessmentResult() error = %v, wantErr %v", err, tt.wantErr)
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

// TODO(lebogg): Convert to table tests
func TestListMetricConfigurations(t *testing.T) {
	var (
		response *orchestrator.ListMetricConfigurationResponse
		err      error
	)
	s := NewService()
	response, err = s.ListMetricConfigurations(context.TODO(), &orchestrator.ListMetricConfigurationRequest{})

	assert.NoError(t, err)
	assert.NotEmpty(t, response.Configurations)
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
		wantErr  bool
	}{
		{
			name: "Store assessment to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:         assessmentResultID1,
						MetricId:   "assessmentResultMetricID",
						EvidenceId: "11111111-1111-1111-1111-111111111111",
						Timestamp:  timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue: toStruct(1.0),
							Operator:    "operator",
							IsDefault:   true,
						},
						NonComplianceComments: "non_compliance_comment",
						Compliant:             true,
						ResourceId:            "resourceID",
					},
				},
			},
			wantErr: false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{
				Status: true,
			},
		},
		{
			name: "Store assessment without metricId to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &assessment.AssessmentResult{
						Id:         assessmentResultID1,
						EvidenceId: "11111111-1111-1111-1111-111111111111",
						Timestamp:  timestamppb.Now(),
						MetricConfiguration: &assessment.MetricConfiguration{
							TargetValue: toStruct(1.0),
							Operator:    "operator",
							IsDefault:   true,
						},
						NonComplianceComments: "non_compliance_comment",
						Compliant:             true,
						ResourceId:            "resourceID",
					},
				},
			},
			wantErr: true,
			wantResp: &orchestrator.StoreAssessmentResultResponse{
				Status:        false,
				StatusMessage: "invalid assessment result: metric id is missing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreAssessmentResult() gotResp = %v, want %v", gotResp, tt.wantResp)
			}

			if err == nil {
				assert.NotNil(t, s.results[assessmentResultID1])
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
					StatusMessage: "invalid assessment result: " + assessment.ErrMetricIdMissing.Error(),
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
				t.Errorf("Got StoreAssessmentResults() error = %v, wantErr %v", err, tt.wantErr)
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
					assert.Equal(t, elem.Status, tt.wantRespMessage[i].Status)
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
				Id:         uuid.NewString(),
				MetricId:   fmt.Sprintf("assessmentResultMetricID-%d", i),
				EvidenceId: "11111111-1111-1111-1111-111111111111",
				Timestamp:  timestamppb.Now(),
				MetricConfiguration: &assessment.MetricConfiguration{
					TargetValue: toStruct(1.0),
					Operator:    fmt.Sprintf("operator%d", i),
					IsDefault:   true,
				},
				NonComplianceComments: "non_compliance_comment",
				Compliant:             true,
				ResourceId:            "resourceID",
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

// CreateCertificateMock creates a mock certificate creation request
func CreateCertificateMock() *orchestrator.Certificate {
	mockHistory := &orchestrator.StateHistory{
		State:         "new",
		TreeId:        "12345",
		Timestamp:     time.Now().String(),
		CertificateID: 1234,
		ID:            12345,
	}

	var mockCertificate = &orchestrator.Certificate{
		Name:          "EUCS",
		CertificateId: "1234",
		ServiceId:     "test service",
		Issuedate:     "2021-11-06",
		Standard:      "EUCS",
		Scope:         "Basic",
		Cab:           "Cab123",
		Description:   "Description",
		StateHistory:  []*orchestrator.StateHistory{mockHistory},
		ID:            1234,
	}

	return mockCertificate
}

func Test_CreateCertificate(t *testing.T) {
	type args struct {
		in0 context.Context
		req *orchestrator.CreateCertificateRequest
	}
	tests := []struct {
		name         string
		args         args
		wantResponse *emptypb.Empty
		wantErr      bool
	}{
		{
			"missing request",
			args{
				context.Background(),
				nil,
			},
			nil,
			true,
		},
		{
			"missing certificate",
			args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{},
			},
			nil,
			true,
		},
		{
			"valid certificate",
			args{
				context.Background(),
				&orchestrator.CreateCertificateRequest{
					Certificate: CreateCertificateMock(),
				},
			},
			&emptypb.Empty{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			gotResponse, err := s.CreateCertificate(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
			CertificateId: "1234",
		},
		CertificateId: "1234",
	})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// 4th case: Certificate updated successfully
	mockCertificate := CreateCertificateMock()
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	// update the certificate's description and send the update request
	mockCertificate.Description = "new description"
	certificate, err = orchestratorService.UpdateCertificate(context.Background(), &orchestrator.UpdateCertificateRequest{
		CertificateId: "1234",
		Certificate:   mockCertificate,
	})
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
	mockCertificate := CreateCertificateMock()
	err = orchestratorService.storage.Create(mockCertificate)
	assert.NoError(t, err)

	// There is a record for certificates in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.NotEmpty(t, listCertificatesResponse.Certificates)

	// Remove record
	_, err = orchestratorService.RemoveCertificate(context.Background(), &orchestrator.RemoveCertificateRequest{CertificateId: mockCertificate.CertificateId})
	assert.NoError(t, err)

	// There is a record for cloud services in the DB (default one)
	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.Empty(t, listCertificatesResponse.Certificates)
}

func Test_GetCertificate(t *testing.T) {
	tests := []struct {
		name string
		req  *orchestrator.GetCertificateRequest
		res  *orchestrator.Certificate
		err  error
	}{
		{
			"invalid request",
			nil,
			nil,
			status.Error(codes.InvalidArgument, orchestrator.ErrRequestIsNil.Error()),
		},
		{
			"certificate not found",
			&orchestrator.GetCertificateRequest{CertificateId: ""},
			nil,
			status.Error(codes.NotFound, "certificate ID is empty"),
		},
		{
			"valid",
			&orchestrator.GetCertificateRequest{CertificateId: "1234"},
			CreateCertificateMock(),
			nil,
		},
	}
	orchestratorService := NewService()
	orchestratorService.storage.Create(CreateCertificateMock())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := orchestratorService.GetCertificate(context.Background(), tt.req)

			if tt.err == nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}

			if tt.res != nil {
				assert.NotEmpty(t, res.CertificateId)
			}

			assert.True(t, proto.Equal(res, tt.res), "%v != %v", res, tt.res)
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
	err = orchestratorService.storage.Create(CreateCertificateMock())
	assert.NoError(t, err)

	listCertificatesResponse, err = orchestratorService.ListCertificates(context.Background(), &orchestrator.ListCertificatesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, listCertificatesResponse.Certificates)
	assert.NotEmpty(t, listCertificatesResponse.Certificates)
	assert.Equal(t, len(listCertificatesResponse.Certificates), 1)
}
