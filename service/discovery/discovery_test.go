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
	service_assessment "clouditor.io/clouditor/service/assessment"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/structpb"
	"net"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	// pre-configuration for mocking evidence store
	const bufSize = 1024 * 1024 * 2
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	assessment.RegisterAssessmentServer(s, service_assessment.NewService())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// make sure, that we are in the clouditor root folder to find the policies
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
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

	type fields struct {
		hasRPCConnection bool
	}

	tests := []struct {
		name string
		fields fields
		wantResp *discovery.StartDiscoveryResponse
		wantErr bool
		wantErrMessage string
	}{
		{
			name: "no error",
			fields: fields{
				hasRPCConnection: true,
			},
			wantResp: &discovery.StartDiscoveryResponse{
				Successful: true,
			},
			wantErr: false,
			wantErrMessage: "",
		},
		{
			name: "No RPC connection",
			fields: fields{
				hasRPCConnection: false,
			},
			wantResp: nil,
			wantErr: true,
			wantErrMessage: codes.Internal.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if tt.fields.hasRPCConnection {
				assert.NoError(t, s.mockAssessmentStream())
			}

			resp, err := s.Start(nil, nil)
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

// Mocking assessment service
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func (s *Service) mockAssessmentStream() error {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client, err := assessment.NewAssessmentClient(conn).AssessEvidences(ctx)
	if err != nil {
		return err
	}
	s.assessmentStream = client
	return nil
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

// toStruct transforms r to a struct and asserts if it was successful
func toStruct(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.NotNil(t, err)
	}

	return
}
