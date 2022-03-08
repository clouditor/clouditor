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
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestMain(m *testing.M) {
	// make sure, that we are in the clouditor root folder to find the policies
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	server, authService, _ := startBufConnServer()
	err = authService.CreateDefaultUser("clouditor", "clouditor")
	if err != nil {
		panic(err)
	}

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
			name: "Create service with options",
			args: args{
				opts: []ServiceOption{
					WithAssessmentAddress("localhost:9091"),
				},
			},
			want: &Service{
				assessmentAddress: "localhost:9091",
				resources:         make(map[string]voc.IsCloudResource),
				configurations:    make(map[discovery.Discoverer]*Configuration),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.opts...)

			// we cannot compare the scheduler, so we need to nil it
			got.scheduler = nil
			tt.want.scheduler = nil

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartDiscovery(t *testing.T) {
	discoveryService := NewService()

	type fields struct {
		assessmentStream    assessment.Assessment_AssessEvidencesClient
		evidenceStoreStream evidence.EvidenceStore_StoreEvidencesClient
		discoverer          discovery.Discoverer
	}

	tests := []struct {
		name          string
		fields        fields
		checkEvidence bool
	}{
		{
			name: "Err in discoverer",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 0}},
		},
		{
			name: "Err in marshaling the resource containing circular dependencies",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 1}},
		},
		{
			name: "No err in discoverer but no evidence stream to assessment",
			fields: fields{
				discoverer: mockDiscoverer{testCase: 2},
			},
		},
		{
			name: "No err in discoverer but no evidence stream to evidence store available",
			fields: fields{
				assessmentStream: &mockAssessmentStream{},
				discoverer:       mockDiscoverer{testCase: 2}},
		},
		{
			name: "No err in discoverer but streaming to assessment fails",
			fields: fields{
				assessmentStream:    &mockAssessmentStream{},
				evidenceStoreStream: mockEvidenceStoreStream{},
				discoverer:          mockDiscoverer{testCase: 2}},
		},
		{
			name: "No err",
			fields: fields{
				assessmentStream:    &mockAssessmentStream{connectionEstablished: true},
				evidenceStoreStream: mockEvidenceStoreStream{},
				discoverer:          mockDiscoverer{testCase: 2}},
			checkEvidence: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discoveryService.assessmentStream = tt.fields.assessmentStream
			discoveryService.StartDiscovery(tt.fields.discoverer)

			// APIs for assessment and evidence store both send the same evidence. Thus, testing one is enough.
			if tt.checkEvidence {
				e := discoveryService.assessmentStream.(*mockAssessmentStream).sentEvidence
				// Check if UUID has been created
				assert.NotEmpty(t, e.Id)
				// Check if cloud resources / properties are there
				assert.NotEmpty(t, e.Resource)
				// Check if ID of mockDiscovery's resource is mapped to resource id of the evidence
				list, _ := tt.fields.discoverer.List()
				assert.Equal(t, string(list[0].GetID()), e.Resource.GetStructValue().AsMap()["id"].(string))
			}
		})
	}

}

func TestQuery(t *testing.T) {
	s := NewService()
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
			numberOfQueriedResources: 1,
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
			assert.Equal(t, tt.numberOfQueriedResources, len(response.Results.Values))
		})
		// If a bad resource was added it will be removed. Otherwise no-op
		delete(s.resources, "MockResourceId")
	}
}

func TestStart(t *testing.T) {

	type envVariable struct {
		hasEnvVariable   bool
		envVariableKey   string
		envVariableValue string
	}

	type fields struct {
		hasRPCConnection bool
		envVariables     []envVariable
	}

	tests := []struct {
		name           string
		fields         fields
		providers      []string
		wantResp       *discovery.StartDiscoveryResponse
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Azure authorizer from file",
			fields: fields{
				hasRPCConnection: true,
				envVariables: []envVariable{
					// We must set AZURE_AUTH_LOCATION to the Azure credentials test file and the set HOME to a wrong path so that the Azure authorizer passes and the K8S authorizer fails
					{
						hasEnvVariable:   true,
						envVariableKey:   "AZURE_AUTH_LOCATION",
						envVariableValue: "service/discovery/testdata/credentials_test_file",
					},
				},
			},
			providers:      []string{ProviderAzure},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "No Azure authorizer",
			fields: fields{
				hasRPCConnection: true,
				envVariables: []envVariable{
					// We must set AZURE_AUTH_LOCATION and HOME to a wrong path so that both Azure authorizer fail
					{
						hasEnvVariable:   true,
						envVariableKey:   "AZURE_AUTH_LOCATION",
						envVariableValue: "",
					},
					{
						hasEnvVariable:   true,
						envVariableKey:   "HOME",
						envVariableValue: "",
					},
				},
			},
			providers:      []string{ProviderAzure},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "could not authenticate to Azure",
		},
		{
			name: "K8s authorizer",
			fields: fields{
				hasRPCConnection: true,
			},
			providers:      []string{ProviderK8S},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "No K8s authorizer",
			fields: fields{
				hasRPCConnection: true,
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
			name: "No RPC connection",
			fields: fields{
				hasRPCConnection: false,
			},
			providers:      []string{"aws", "azure", "k8s"},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "could not initialize stream to Assessment",
		},
		{
			name: "Request with 2 providers",
			fields: fields{
				hasRPCConnection: true,
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
			name: "Empty request",
			fields: fields{
				hasRPCConnection: true,
			},
			providers:      []string{},
			wantResp:       &discovery.StartDiscoveryResponse{Successful: true},
			wantErr:        false,
			wantErrMessage: "",
		},
		{
			name: "Request with wrong provider name",
			fields: fields{
				hasRPCConnection: true,
			},
			providers:      []string{"falseProvider"},
			wantResp:       nil,
			wantErr:        true,
			wantErrMessage: "provider falseProvider not known",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(WithProviders(tt.providers))
			if tt.fields.hasRPCConnection {
				assert.NoError(t, s.initAssessmentStream(grpc.WithContextDialer(bufConnDialer)))
			}

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

func TestService_initAssessmentStream(t *testing.T) {
	type fields struct {
		hasRPCConnection bool
		username         string
		password         string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "No RPC connection",
			fields: fields{
				hasRPCConnection: false,
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				// We are looking for a connection refused message
				innerErr := errors.Unwrap(err)
				s, _ := status.FromError(innerErr)

				if s.Code() != codes.Unavailable {
					tt.Errorf("Status should be codes.Unavailable: %v", s.Code())
					return false
				}

				return true
			},
		},
		{
			name: "Authenticated RPC connection with valid user",
			fields: fields{
				hasRPCConnection: true,
				username:         "clouditor",
				password:         "clouditor",
			},
			wantErr: nil,
		},
		{
			name: "Authenticated RPC connection with invalid user",
			fields: fields{
				hasRPCConnection: true,
				username:         "notclouditor",
				password:         "password",
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(tt, err.Error(), "could not set up stream for assessing evidences: rpc error: code = Unauthenticated desc = transport: per-RPC creds failed due to error: error while logging in: rpc error: code = Unauthenticated desc = login failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(WithInternalAuthorizer(
				"bufnet",
				tt.fields.username,
				tt.fields.password,
				grpc.WithContextDialer(bufConnDialer),
			))

			var opts []grpc.DialOption
			if tt.fields.hasRPCConnection {
				// Make this a valid RPC connection by connecting to our bufnet service
				opts = []grpc.DialOption{grpc.WithContextDialer(bufConnDialer)}
				s.assessmentAddress = "bufnet"
			}

			err := s.initAssessmentStream(opts...)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}
		})
	}
}

func TestShutdown(t *testing.T) {
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
					CloudResource: &voc.CloudResource{
						ID:   "some-id",
						Name: "some-name",
						Type: []string{"ObjectStorage", "Storage", "Resource"},
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
	sentEvidence *evidence.Evidence
	// We add connectionEstablished to differentiate between the case where evidences can be sent and not
	connectionEstablished bool
	counter               int
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
	e := req.Evidence
	if m.connectionEstablished {
		m.sentEvidence = e
	} else {
		err = fmt.Errorf("mock send error")
	}
	return
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

func (*mockAssessmentStream) SendMsg(_ interface{}) error {
	return nil
}

func (*mockAssessmentStream) RecvMsg(_ interface{}) error {
	return nil
}

// mockEvidenceStoreStream implements EvidenceStore_StoreEvidencesClient interface
type mockEvidenceStoreStream struct {
}

func (mockEvidenceStoreStream) Recv() (*evidence.StoreEvidenceResponse, error) {
	return nil, nil
}

func (mockEvidenceStoreStream) Send(_ *evidence.StoreEvidenceRequest) error {
	return fmt.Errorf("mock send error")
}

func (mockEvidenceStoreStream) CloseAndRecv() (*emptypb.Empty, error) {
	return nil, nil
}

func (mockEvidenceStoreStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (mockEvidenceStoreStream) Trailer() metadata.MD {
	return nil
}

func (mockEvidenceStoreStream) CloseSend() error {
	return nil
}

func (mockEvidenceStoreStream) Context() context.Context {
	return nil
}

func (mockEvidenceStoreStream) SendMsg(_ interface{}) error {
	return nil
}

func (mockEvidenceStoreStream) RecvMsg(_ interface{}) error {
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
