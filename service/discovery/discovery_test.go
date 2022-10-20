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

package discovery

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMain(m *testing.M) {
	clitest.AutoChdir()

	server, _ := startBufConnServer()

	code := m.Run()

	server.Stop()
	os.Exit(code)
}

func TestNewService(t *testing.T) {
	type args struct {
		opts []ServiceOption
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Create service with option 'WithAssessmentAddress'",
			args: args{
				opts: []ServiceOption{
					WithAssessmentAddress("localhost:9091"),
				},
			},
			want: &Service{
				assessmentAddress: grpcTarget{target: "localhost:9091"},
				resources:         make(map[string]voc.IsCloudResource),
				configurations:    make(map[discovery.Discoverer]*Configuration),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// we cannot compare the scheduler, so we first check if it is not empty and then nil it
			assert.NotEmpty(t, got.scheduler)
			got.scheduler = nil
			tt.want.scheduler = nil

			// we cannot compare the assessment streams, so we first check if it is not empty and then nil it
			assert.NotEmpty(t, got.assessmentStreams)
			got.assessmentStreams = nil
			tt.want.assessmentStreams = nil

			// we cannot compare the Events, so we first check if it is not empty and then nil it
			assert.Empty(t, got.Events)
			got.Events = nil
			tt.want.Events = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_StartDiscovery(t *testing.T) {
	type fields struct {
		discoverer discovery.Discoverer
	}

	tests := []struct {
		name          string
		fields        fields
		checkEvidence bool
	}{
		{
			name: "Err in discoverer",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 0},
			},
		},
		{
			name: "Err in marshaling the resource containing circular dependencies",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 1},
			},
		},
		{
			name: "No err",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 2},
			},
			checkEvidence: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
			mockStream.Prepare()

			svc := NewService()
			svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
			_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
				return mockStream, nil
			})
			svc.assessmentAddress = grpcTarget{target: "mock"}
			go svc.StartDiscovery(tt.fields.discoverer)

			if tt.checkEvidence {
				mockStream.Wait()

				want, _ := tt.fields.discoverer.List()

				got := mockStream.sentEvidences
				assert.Equal(t, len(want), len(got))

				// Retrieve the last one
				e := got[len(got)-1]

				// Check, if evidence was sent
				assert.NotNil(t, e)
				// Check if UUID has been created
				assert.NotEmpty(t, e.Id)
				// Check if cloud resources / properties are there
				assert.NotEmpty(t, e.Resource)
				// Check if ID of mockDiscovery's resource is mapped to resource id of the evidence

				// Only the last element sent can be checked
				assert.Equal(t, string(want[len(want)-1].GetID()), e.Resource.GetStructValue().AsMap()["id"].(string))
			}
		})
	}
}

func TestService_Query(t *testing.T) {
	s := NewService(WithAssessmentAddress("bufnet", grpc.WithContextDialer(bufConnDialer)))
	s.StartDiscovery(mockDiscoverer{testCase: 2})

	type fields struct {
		resources    map[string]voc.IsCloudResource
		queryRequest *discovery.QueryRequest
	}
	tests := []struct {
		name string
		fields
		numberOfQueriedResources int
		wantErr                  bool
	}{
		{
			name: "Err when unmarshalling",
			fields: fields{
				queryRequest: &discovery.QueryRequest{},
				resources: map[string]voc.IsCloudResource{
					"MockResourceId": wrongFormattedResource(),
				},
			},
			numberOfQueriedResources: 1,
			wantErr:                  true,
		},
		{
			name:                     "Filter type",
			fields:                   fields{queryRequest: &discovery.QueryRequest{FilteredType: "Compute"}},
			numberOfQueriedResources: 0,
		},
		{
			name:                     "No filtering",
			fields:                   fields{queryRequest: &discovery.QueryRequest{}},
			numberOfQueriedResources: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add resources (bad formatted ones)
			if res := tt.fields.resources; res != nil {
				for k, v := range res {
					s.resources[k] = v
				}
			}
			response, err := s.Query(context.TODO(), &discovery.QueryRequest{FilteredType: tt.fields.queryRequest.FilteredType})
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}
			assert.Equal(t, tt.numberOfQueriedResources, len(response.Results))
		})
		// If a bad resource was added it will be removed. Otherwise no-op
		delete(s.resources, "MockResourceId")
	}
}

func TestService_Start(t *testing.T) {

	type envVariable struct {
		hasEnvVariable   bool
		envVariableKey   string
		envVariableValue string
	}

	type fields struct {
		envVariables []envVariable
	}

	tests := []struct {
		name           string
		fields         fields
		req            *discovery.StartDiscoveryRequest
		providers      []string
		wantResp       *discovery.StartDiscoveryResponse
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:           "No Azure authorizer",
			providers:      []string{ProviderAzure},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "could not authenticate to Azure:",
		},
		{
			name: "Azure authorizer from ENV",
			fields: fields{
				envVariables: []envVariable{
					// We must set AZURE_AUTH_LOCATION to the Azure credentials test file and the set HOME to a
					// wrong path so that the Azure authorizer passes and the K8S authorizer fails
					{
						hasEnvVariable:   true,
						envVariableKey:   "AZURE_TENANT_ID",
						envVariableValue: "tenant-id-123",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "AZURE_CLIENT_ID",
						envVariableValue: "client-id-123",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "AZURE_CLIENT_SECRET",
						envVariableValue: "client-secret-456",
					},
				},
			},
			providers:      []string{ProviderAzure},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "No K8s authorizer",
			fields: fields{
				envVariables: []envVariable{
					// We must set HOME to a wrong path so that the K8S authorizer fails
					{
						hasEnvVariable:   true,
						envVariableKey:   "HOME",
						envVariableValue: "",
					},
				},
			},
			providers:      []string{ProviderK8S},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "could not authenticate to Kubernetes",
		},
		{
			name: "Request with 2 providers",
			fields: fields{
				envVariables: []envVariable{
					// We must set HOME to a wrong path so that the AWS and k8s authorizer fails in all systems, regardless if AWS and k8s paths are set or not
					{
						hasEnvVariable:   true,
						envVariableKey:   "HOME",
						envVariableValue: "",
					},
				},
			},
			providers:      []string{"aws", "k8s"},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "could not authenticate to",
		},
		{
			name:           "Empty request",
			fields:         fields{},
			providers:      []string{},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name:           "Request with wrong provider name",
			fields:         fields{},
			providers:      []string{"falseProvider"},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "provider falseProvider not known",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(WithProviders(tt.providers))

			for _, env := range tt.fields.envVariables {
				if env.hasEnvVariable {
					t.Setenv(env.envVariableKey, env.envVariableValue)
				}
			}

			resp, err := s.Start(context.TODO(), nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Got Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				assert.Equal(t, tt.wantResp, resp)
			}

			if err != nil {
				assert.Contains(t, err.Error(), tt.wantErrMessage)
			}
		})
	}
}

func TestService_Shutdown(t *testing.T) {
	service := NewService()
	service.Shutdown()

	assert.False(t, service.scheduler.IsRunning())

}

// TestHandleError tests handleError (if all error cases are executed/printed)
func TestHandleError(t *testing.T) {
	type args struct {
		err  error
		dest string
	}
	tests := []struct {
		name           string
		args           args
		wantErrSnippet string
	}{
		{
			name: "handleInternalError",
			args: args{
				err:  status.Error(codes.Internal, "internal error"),
				dest: "SomeDestination",
			},
			wantErrSnippet: "internal",
		},
		{
			name: "handleInvalidError",
			args: args{
				err:  status.Errorf(codes.InvalidArgument, "invalid argument"),
				dest: "SomeDestination",
			},
			wantErrSnippet: "invalid",
		},
		{
			name: "handleSomeOtherErr",
			args: args{
				err:  errors.New("some other error"),
				dest: "SomeDestination",
			},
			wantErrSnippet: "some other error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleError(tt.args.err, tt.args.dest)
			assert.Contains(t, err.Error(), tt.wantErrSnippet)
		})
	}
}

// mockDiscoverer implements Discoverer and mocks the API to cloud resources
type mockDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	testCase int
}

func (mockDiscoverer) Name() string { return "just mocking" }

func (m mockDiscoverer) List() ([]voc.IsCloudResource, error) {
	switch m.testCase {
	case 0:
		return nil, fmt.Errorf("mock error in List()")
	case 1:
		return []voc.IsCloudResource{wrongFormattedResource()}, nil
	case 2:
		return []voc.IsCloudResource{
			&voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:   "some-id",
						Name: "some-name",
						Type: []string{"ObjectStorage", "Storage", "Resource"},
					},
				},
			},
			&voc.StorageService{
				Storages: []voc.ResourceID{"some-id"},
				NetworkService: &voc.NetworkService{
					Networking: &voc.Networking{
						Resource: &voc.Resource{
							ID:   "some-storage-account-id",
							Name: "some-storage-account-name",
							Type: []string{"StorageService", "NetworkService", "Networking", "Resource"},
						},
					},
				},
				HttpEndpoint: &voc.HttpEndpoint{
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   false,
						Enabled:    true,
						TlsVersion: "TLS1_2",
					},
				},
			},
		}, nil
	default:
		return nil, nil
	}
}

func wrongFormattedResource() voc.IsCloudResource {
	res1 := mockIsCloudResource{Another: nil}
	res2 := mockIsCloudResource{Another: &res1}
	res1.Another = &res2
	return res1
}

// mockAssessmentStream implements Assessment_AssessEvidencesClient interface
type mockAssessmentStream struct {
	// We add sentEvidence field to test the evidence that would be sent over gRPC
	sentEvidences []*evidence.Evidence
	// We add connectionEstablished to differentiate between the case where evidences can be sent and not
	connectionEstablished bool
	counter               int
	expected              int
	wg                    sync.WaitGroup
}

func (m *mockAssessmentStream) Prepare() {
	m.wg.Add(m.expected)
}

func (m *mockAssessmentStream) Wait() {
	m.wg.Wait()
}

func (m *mockAssessmentStream) Recv() (*assessment.AssessEvidenceResponse, error) {
	if m.counter == 0 {
		m.counter++
		return &assessment.AssessEvidenceResponse{
			Status:        assessment.AssessEvidenceResponse_FAILED,
			StatusMessage: "mockError1",
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &assessment.AssessEvidenceResponse{
			Status: assessment.AssessEvidenceResponse_ASSESSED,
		}, nil
	} else {
		return nil, io.EOF
	}
}

func (m *mockAssessmentStream) Send(req *assessment.AssessEvidenceRequest) (err error) {
	return m.SendMsg(req)
}

func (*mockAssessmentStream) CloseAndRecv() (*emptypb.Empty, error) {
	return nil, nil
}

func (*mockAssessmentStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (*mockAssessmentStream) Trailer() metadata.MD {
	return nil
}

func (*mockAssessmentStream) CloseSend() error {
	return nil
}

func (*mockAssessmentStream) Context() context.Context {
	return nil
}

func (m *mockAssessmentStream) SendMsg(req interface{}) (err error) {
	e := req.(*assessment.AssessEvidenceRequest).Evidence
	if m.connectionEstablished {
		m.sentEvidences = append(m.sentEvidences, e)
	} else {
		err = fmt.Errorf("mock send error")
	}

	m.wg.Done()

	return
}

func (*mockAssessmentStream) RecvMsg(_ interface{}) error {
	return nil
}

// mockIsCloudResource implements mockIsCloudResource interface.
// It is used for json.marshal to fail since it contains circular dependency
type mockIsCloudResource struct {
	Another *mockIsCloudResource `json:"Another"`
}

func (mockIsCloudResource) GetID() voc.ResourceID {
	return "MockResourceId"
}

func (mockIsCloudResource) GetServiceID() string {
	return "MockServiceId"
}

func (mockIsCloudResource) GetName() string {
	return ""
}

func (mockIsCloudResource) GetType() []string {
	return nil
}

func (mockIsCloudResource) HasType(_ string) bool {
	return false
}

func (mockIsCloudResource) GetCreationTime() *time.Time {
	return nil
}

func (mockIsCloudResource) Related() []string {
	return []string{}
}
