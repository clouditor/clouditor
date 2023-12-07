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
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/testutil"
	"clouditor.io/clouditor/internal/testutil/clitest"
	"clouditor.io/clouditor/internal/testutil/servicetest"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/persistence"
	"clouditor.io/clouditor/service"
	"clouditor.io/clouditor/voc"

	"github.com/go-co-op/gocron"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
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
		want assert.ValueAssertionFunc
	}{
		{
			name: "Create service with option 'WithAssessmentAddress'",
			args: args{
				opts: []ServiceOption{
					WithAssessmentAddress("localhost:9091"),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, "localhost:9091", s.assessment.Target)
			},
		},
		{
			name: "Create service with option 'WithDefaultCloudServiceID'",
			args: args{
				opts: []ServiceOption{
					WithCloudServiceID(testdata.MockCloudServiceID1),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, testdata.MockCloudServiceID1, s.csID)
			},
		},
		{
			name: "Create service with option 'WithAuthorizationStrategy'",
			args: args{
				opts: []ServiceOption{
					WithAuthorizationStrategy(&service.AuthorizationStrategyJWT{AllowAllKey: "test"}),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, &service.AuthorizationStrategyJWT{AllowAllKey: "test"}, s.authz)
			},
		},
		{
			name: "Create service with option 'WithStorage'",
			args: args{
				opts: []ServiceOption{
					WithStorage(testutil.NewInMemoryStorage(t)),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.NotNil(t, s.storage)
			},
		},
		{
			name: "Create service with option 'WithDiscoveryInterval'",
			args: args{
				opts: []ServiceOption{
					WithDiscoveryInterval(time.Duration(8)),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				s := i1.(*Service)
				return assert.Equal(t, s.discoveryInterval, time.Duration(8))
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

func TestService_StartDiscovery(t *testing.T) {
	type fields struct {
		discoverer discovery.Discoverer
		csID       string
	}

	tests := []struct {
		name          string
		fields        fields
		checkEvidence bool
	}{
		{
			name: "Err in discoverer",
			fields: fields{
				discoverer: &mockDiscoverer{testCase: 0, csID: discovery.DefaultCloudServiceID},
				csID:       discovery.DefaultCloudServiceID,
			},
		},
		{
			name: "Err in marshaling the resource containing circular dependencies",
			fields: fields{
				discoverer: &mockDiscoverer{testCase: 1, csID: discovery.DefaultCloudServiceID},
				csID:       discovery.DefaultCloudServiceID,
			},
		},
		{
			name: "No err with default cloud service ID",
			fields: fields{
				discoverer: &mockDiscoverer{testCase: 2, csID: discovery.DefaultCloudServiceID},
				csID:       discovery.DefaultCloudServiceID,
			},
			checkEvidence: true,
		},
		{
			name: "No err with custom cloud service ID",
			fields: fields{
				discoverer: &mockDiscoverer{testCase: 2, csID: testdata.MockCloudServiceID1},
				csID:       testdata.MockCloudServiceID1,
			},
			checkEvidence: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStream := &mockAssessmentStream{connectionEstablished: true, expected: 2}
			mockStream.Prepare()

			svc := NewService()
			svc.csID = tt.fields.csID
			svc.assessmentStreams = api.NewStreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]()
			_, _ = svc.assessmentStreams.GetStream("mock", "Assessment", func(target string, additionalOpts ...grpc.DialOption) (stream assessment.Assessment_AssessEvidencesClient, err error) {
				return mockStream, nil
			})
			svc.assessment = &api.RPCConnection[assessment.AssessmentClient]{Target: "mock"}
			go svc.StartDiscovery(tt.fields.discoverer)

			if tt.checkEvidence {
				mockStream.Wait()
				want, _ := tt.fields.discoverer.List()

				got := mockStream.sentEvidences
				assert.Equal(t, len(want), len(got))

				// Retrieve the last one
				eWant := want[len(want)-1]
				eGot := got[len(got)-1]
				err := eGot.Validate()
				assert.NotNil(t, eGot)
				assert.NoError(t, err)

				// Only the last element sent can be checked
				assert.Equal(t, string(eWant.GetID()), eGot.Resource.GetStructValue().AsMap()["id"].(string))

				// Assert cloud service ID
				assert.Equal(t, tt.fields.csID, eGot.CloudServiceId)
				assert.Equal(t, tt.fields.csID, eGot.Resource.GetStructValue().AsMap()["serviceId"].(string))
			}
		})
	}
}

func TestService_ListResources(t *testing.T) {
	type fields struct {
		authz service.AuthorizationStrategy
		csID  string
	}
	type args struct {
		req *discovery.ListResourcesRequest
	}
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		numberOfQueriedResources int
		wantErr                  assert.ErrorAssertionFunc
	}{
		{
			name: "Filter type, allow all",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				csID:  testdata.MockCloudServiceID1,
			},
			args: args{req: &discovery.ListResourcesRequest{
				Filter: &discovery.ListResourcesRequest_Filter{
					Type: util.Ref("Storage"),
				},
			}},
			numberOfQueriedResources: 1,
			wantErr:                  assert.NoError,
		},
		{
			name: "Filter cloud service, allow",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
				csID:  testdata.MockCloudServiceID1,
			},
			args: args{req: &discovery.ListResourcesRequest{
				Filter: &discovery.ListResourcesRequest_Filter{
					CloudServiceId: util.Ref(testdata.MockCloudServiceID1),
				},
			}},
			numberOfQueriedResources: 2,
			wantErr:                  assert.NoError,
		},
		{
			name: "Filter cloud service, not allowed",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID1),
				csID:  testdata.MockCloudServiceID1,
			},
			args: args{req: &discovery.ListResourcesRequest{
				Filter: &discovery.ListResourcesRequest_Filter{
					CloudServiceId: util.Ref(testdata.MockCloudServiceID2),
				},
			}},
			numberOfQueriedResources: 0,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, service.ErrPermissionDenied)
			},
		},
		{
			name: "No filtering, allow all",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				csID:  testdata.MockCloudServiceID1,
			},
			args:                     args{req: &discovery.ListResourcesRequest{}},
			numberOfQueriedResources: 2,
			wantErr:                  assert.NoError,
		},
		{
			name: "No filtering, allow different cloud service, empty result",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
				csID:  testdata.MockCloudServiceID1,
			},
			args:                     args{req: &discovery.ListResourcesRequest{}},
			numberOfQueriedResources: 0,
			wantErr:                  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(WithAssessmentAddress("bufnet", grpc.WithContextDialer(bufConnDialer)))
			s.authz = tt.fields.authz
			s.csID = tt.fields.csID
			s.StartDiscovery(&mockDiscoverer{testCase: 2, csID: tt.fields.csID})

			response, err := s.ListResources(context.TODO(), tt.args.req)
			tt.wantErr(t, err)

			if err == nil {
				assert.Equal(t, tt.numberOfQueriedResources, len(response.Results))
			}
		})
	}
}

func TestService_Shutdown(t *testing.T) {
	service := NewService()
	service.Shutdown()

	assert.False(t, service.scheduler.IsRunning())

}

// mockDiscoverer implements Discoverer and mocks the API to cloud resources
type mockDiscoverer struct {
	// testCase allows for different implementations for table tests in TestStartDiscovery
	testCase int
	csID     string
}

func (*mockDiscoverer) Name() string { return "just mocking" }

func (m *mockDiscoverer) List() ([]voc.IsCloudResource, error) {
	switch m.testCase {
	case 0:
		return nil, fmt.Errorf("mock error in List()")
	case 1:
		return []voc.IsCloudResource{wrongFormattedResource()}, nil
	case 2:
		return []voc.IsCloudResource{
			&voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: discovery.NewResource(m, "some-id", "some-name", nil, voc.GeoLocation{}, nil, "", []string{"ObjectStorage", "Storage", "Resource"}, map[string][]interface{}{"raw": {"raw"}}),
				},
			},
			&voc.ObjectStorageService{
				StorageService: &voc.StorageService{
					Storage: []voc.ResourceID{"some-id"},
					NetworkService: &voc.NetworkService{
						Networking: &voc.Networking{
							Resource: discovery.NewResource(m, "some-storage-account-id", "some-storage-account-name", nil, voc.GeoLocation{}, nil, "", []string{"StorageService", "NetworkService", "Networking", "Resource"}, map[string][]interface{}{"raw": {"raw"}}),
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

// CloudServiceID is an implementation for discovery.Discoverer
func (d *mockDiscoverer) CloudServiceID() string {
	return d.csID
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

func (m *mockAssessmentStream) Recv() (*assessment.AssessEvidencesResponse, error) {
	if m.counter == 0 {
		m.counter++
		return &assessment.AssessEvidencesResponse{
			Status:        assessment.AssessEvidencesResponse_FAILED,
			StatusMessage: "mockError1",
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &assessment.AssessEvidencesResponse{
			Status: assessment.AssessEvidencesResponse_ASSESSED,
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

func (mockIsCloudResource) SetServiceID(_ string) {

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

func (mockIsCloudResource) GetRaw() string {
	return ""
}

func (mockIsCloudResource) SetRaw(_ string) {
}

func (mockIsCloudResource) Related() []string {
	return []string{}
}

func Test_toDiscoveryResource(t *testing.T) {
	type args struct {
		resource voc.IsCloudResource
	}
	tests := []struct {
		name    string
		args    args
		wantR   *discovery.Resource
		wantV   *structpb.Value
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				resource: &voc.VirtualMachine{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:   "my-resource-id",
							Name: "my-resource-name",
							Type: voc.VirtualMachineType,
						},
					},
				},
			},
			wantV: &structpb.Value{
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"activityLogging":  structpb.NewNullValue(),
							"automaticUpdates": structpb.NewNullValue(),
							"blockStorage":     structpb.NewNullValue(),
							"bootLogging":      structpb.NewNullValue(),
							"creationTime":     structpb.NewNumberValue(0),
							"geoLocation": structpb.NewStructValue(&structpb.Struct{
								Fields: map[string]*structpb.Value{
									"region": structpb.NewStringValue(""),
								}}),
							"id":                structpb.NewStringValue("my-resource-id"),
							"labels":            structpb.NewNullValue(),
							"malwareProtection": structpb.NewNullValue(),
							"name":              structpb.NewStringValue("my-resource-name"),
							"networkInterfaces": structpb.NewNullValue(),
							"osLogging":         structpb.NewNullValue(),
							"parent":            structpb.NewStringValue(""),
							"raw":               structpb.NewStringValue(""),
							"resourceLogging":   structpb.NewNullValue(),
							"serviceId":         structpb.NewStringValue(""),
							"type": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
								structpb.NewStringValue("VirtualMachine"),
								structpb.NewStringValue("Compute"),
								structpb.NewStringValue("Resource"),
							}}),
						},
					},
				},
			},
			wantR: &discovery.Resource{
				Id:           "my-resource-id",
				ResourceType: "VirtualMachine,Compute,Resource",
				Properties: &structpb.Value{
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"activityLogging":  structpb.NewNullValue(),
								"automaticUpdates": structpb.NewNullValue(),
								"blockStorage":     structpb.NewNullValue(),
								"bootLogging":      structpb.NewNullValue(),
								"creationTime":     structpb.NewNumberValue(0),
								"geoLocation": structpb.NewStructValue(&structpb.Struct{
									Fields: map[string]*structpb.Value{
										"region": structpb.NewStringValue(""),
									}}),
								"id":                structpb.NewStringValue("my-resource-id"),
								"labels":            structpb.NewNullValue(),
								"malwareProtection": structpb.NewNullValue(),
								"name":              structpb.NewStringValue("my-resource-name"),
								"networkInterfaces": structpb.NewNullValue(),
								"osLogging":         structpb.NewNullValue(),
								"parent":            structpb.NewStringValue(""),
								"raw":               structpb.NewStringValue(""),
								"resourceLogging":   structpb.NewNullValue(),
								"serviceId":         structpb.NewStringValue(""),
								"type": structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
									structpb.NewStringValue("VirtualMachine"),
									structpb.NewStringValue("Compute"),
									structpb.NewStringValue("Resource"),
								}}),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotV, err := toDiscoveryResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("toDiscoveryResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !proto.Equal(gotR, tt.wantR) {
				t.Errorf("toDiscoveryResource() r = %v, want %v", gotR, tt.wantR)
			}

			if !proto.Equal(gotV, tt.wantV) {
				t.Errorf("toDiscoveryResource() v = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestService_Start(t *testing.T) {
	type envVariable struct {
		hasEnvVariable   bool
		envVariableKey   string
		envVariableValue string
	}
	type fields struct {
		assessmentStreams *api.StreamsOf[assessment.Assessment_AssessEvidencesClient, *assessment.AssessEvidenceRequest]
		assessment        *api.RPCConnection[assessment.AssessmentClient]
		storage           persistence.Storage
		scheduler         *gocron.Scheduler
		authz             service.AuthorizationStrategy
		providers         []string
		discoveryInterval time.Duration
		Events            chan *DiscoveryEvent
		csID              string
		envVariables      []envVariable
	}
	type args struct {
		ctx context.Context
		req *discovery.StartDiscoveryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO(all): How to test for Azure and AWS authorizer failures and K8S authorizer without failure?
		{
			name: "Invalid request",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
			},
			args: args{
				ctx: context.Background(),
				req: nil,
			},
			want: assert.Nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Request with wrong provider name",
			fields: fields{
				authz:     servicetest.NewAuthorizationStrategy(true),
				scheduler: gocron.NewScheduler(time.UTC),
				providers: []string{"falseProvider"},
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "provider falseProvider not known")
			},
		},
		{
			name: "Wrong permission",
			fields: fields{
				authz:     servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
				scheduler: gocron.NewScheduler(time.UTC),
				providers: []string{},
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "access denied")
			},
		},
		{
			name: "discovery interval error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				scheduler:         gocron.NewScheduler(time.UTC),
				providers:         []string{ProviderAzure},
				discoveryInterval: time.Duration(-5 * time.Minute),
				envVariables: []envVariable{
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
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not schedule job for ", ".Every() interval must be greater than 0")
			},
		},
		{
			name: "K8S authorizer error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				scheduler:         gocron.NewScheduler(time.UTC),
				providers:         []string{ProviderK8S},
				discoveryInterval: time.Duration(5 * time.Minute),
				envVariables: []envVariable{
					// We must set HOME to a wrong path so that the K8S authorizer fails
					{
						hasEnvVariable:   true,
						envVariableKey:   "HOME",
						envVariableValue: "",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not authenticate to Kubernetes")
			},
		},
		{
			name: "Happy path: no discovery interval error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				scheduler:         gocron.NewScheduler(time.UTC),
				providers:         []string{ProviderAzure},
				discoveryInterval: time.Duration(5 * time.Minute),
				envVariables: []envVariable{
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
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				resp, ok := i1.(*discovery.StartDiscoveryResponse)
				assert.True(t, ok)
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, resp)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: Azure authorizer from ENV",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				scheduler:         gocron.NewScheduler(time.UTC),
				providers:         []string{ProviderAzure},
				discoveryInterval: time.Duration(5 * time.Minute),
				envVariables: []envVariable{
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
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				resp, ok := i1.(*discovery.StartDiscoveryResponse)
				assert.True(t, ok)
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, resp)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: Azure with resource group",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				scheduler:         gocron.NewScheduler(time.UTC),
				providers:         []string{ProviderAzure},
				discoveryInterval: time.Duration(5 * time.Minute),
				envVariables: []envVariable{
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
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{
					ResourceGroup: util.Ref("testResourceGroup"),
				},
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				resp, ok := i1.(*discovery.StartDiscoveryResponse)
				assert.True(t, ok)
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, resp)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				scheduler:         tt.fields.scheduler,
				authz:             tt.fields.authz,
				providers:         tt.fields.providers,
				discoveryInterval: tt.fields.discoveryInterval,
				Events:            tt.fields.Events,
				csID:              tt.fields.csID,
			}

			// Set env variables
			for _, env := range tt.fields.envVariables {
				if env.hasEnvVariable {
					t.Setenv(env.envVariableKey, env.envVariableValue)
				}
			}

			gotResp, err := svc.Start(tt.args.ctx, tt.args.req)

			tt.want(t, gotResp)
			tt.wantErr(t, err)
		})
	}
}
