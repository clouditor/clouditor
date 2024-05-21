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

	"clouditor.io/clouditor/v2/api"
	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/clitest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/launcher"
	"clouditor.io/clouditor/v2/persistence"
	"clouditor.io/clouditor/v2/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
		opts []service.Option[*Service]
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*Service]
	}{
		{
			name: "Create service with option 'WithAssessmentAddress'",
			args: args{
				opts: []service.Option[*Service]{
					WithAssessmentAddress("localhost:9091"),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, "localhost:9091", got.assessment.Target)
			},
		},
		{
			name: "Create service with option 'WithDefaultCloudServiceID'",
			args: args{
				opts: []service.Option[*Service]{
					WithCloudServiceID(testdata.MockCloudServiceID1),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, testdata.MockCloudServiceID1, got.csID)
			},
		},
		{
			name: "Create service with option 'WithAuthorizationStrategy'",
			args: args{
				opts: []service.Option[*Service]{
					WithAuthorizationStrategy(&service.AuthorizationStrategyJWT{AllowAllKey: "test"}),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal[service.AuthorizationStrategy](t, &service.AuthorizationStrategyJWT{AllowAllKey: "test"}, got.authz)
			},
		},
		{
			name: "Create service with option 'WithProviders' and one provider given",
			args: args{
				opts: []service.Option[*Service]{
					WithProviders([]string{"azure"}),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, []string{"azure"}, got.providers)
			},
		},
		{
			name: "Create service with option 'WithProviders' and no provider given",
			args: args{
				opts: []service.Option[*Service]{
					WithProviders([]string{}),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, []string{}, got.providers)
			},
		},
		{
			name: "Create service with option 'WithStorage'",
			args: args{
				opts: []service.Option[*Service]{
					WithStorage(testutil.NewInMemoryStorage(t)),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.NotNil(t, got.storage)
			},
		},
		{
			name: "Create service with option 'WithAdditionalDiscoverers'",
			args: args{
				opts: []service.Option[*Service]{
					WithAdditionalDiscoverers([]discovery.Discoverer{&discoverytest.TestDiscoverer{ServiceId: config.DefaultCloudServiceID}}),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Contains(t, got.discoverers, &discoverytest.TestDiscoverer{ServiceId: config.DefaultCloudServiceID})
			},
		},
		{
			name: "Create service with option 'WithDiscoveryInterval'",
			args: args{
				opts: []service.Option[*Service]{
					WithDiscoveryInterval(time.Duration(8)),
				},
			},
			want: func(t *testing.T, got *Service) bool {
				return assert.Equal(t, time.Duration(8), got.discoveryInterval)
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
				discoverer: &discoverytest.TestDiscoverer{TestCase: 0, ServiceId: config.DefaultCloudServiceID},
				csID:       config.DefaultCloudServiceID,
			},
		},
		{
			name: "No err with default cloud service ID",
			fields: fields{
				discoverer: &discoverytest.TestDiscoverer{TestCase: 2, ServiceId: config.DefaultCloudServiceID},
				csID:       config.DefaultCloudServiceID,
			},
			checkEvidence: true,
		},
		{
			name: "No err with custom cloud service ID",
			fields: fields{
				discoverer: &discoverytest.TestDiscoverer{TestCase: 2, ServiceId: testdata.MockCloudServiceID1},
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
				err := api.Validate(eGot)
				assert.NotNil(t, eGot)
				assert.NoError(t, err)

				m, err := eGot.Resource.UnmarshalNew()
				assert.NoError(t, err)
				or := m.(ontology.IsResource)

				// Only the last element sent can be checked
				assert.Equal(t, string(eWant.GetId()), or.GetId())

				// Assert cloud service ID
				assert.Equal(t, tt.fields.csID, eGot.CloudServiceId)
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
		wantErr                  assert.WantErr
	}{
		{
			name: "Filter type, allow all",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(true),
				csID:  testdata.MockCloudServiceID1,
			},
			args: args{req: &discovery.ListResourcesRequest{
				Filter: &discovery.ListResourcesRequest_Filter{
					// TODO(oxisto): This is a problem now, since we are only persisting the leaf node type, so we cannot "see" the inherited resource types anymore
					Type: util.Ref("Storage"),
				},
			}},
			numberOfQueriedResources: 1,
			wantErr:                  assert.Nil[error],
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
			wantErr:                  assert.Nil[error],
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
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorIs(t, gotErr, service.ErrPermissionDenied)
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
			wantErr:                  assert.Nil[error],
		},
		{
			name: "No filtering, allow different cloud service, empty result",
			fields: fields{
				authz: servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
				csID:  testdata.MockCloudServiceID1,
			},
			args:                     args{req: &discovery.ListResourcesRequest{}},
			numberOfQueriedResources: 0,
			wantErr:                  assert.Nil[error],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(WithAssessmentAddress(testdata.MockGRPCTarget, grpc.WithContextDialer(bufConnDialer)))
			s.authz = tt.fields.authz
			s.csID = tt.fields.csID
			s.StartDiscovery(&discoverytest.TestDiscoverer{TestCase: 2, ServiceId: tt.fields.csID})

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

	assert.Empty(t, service.tickers)
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
			Status:        assessment.AssessmentStatus_ASSESSMENT_STATUS_FAILED,
			StatusMessage: "mockError1",
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &assessment.AssessEvidencesResponse{
			Status: assessment.AssessmentStatus_ASSESSMENT_STATUS_ASSESSED,
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
		tickers           map[discovery.Discoverer]*discoveryTicker
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
		want    assert.Want[*discovery.StartDiscoveryResponse]
		wantErr assert.WantErr
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
			want: assert.Nil[*discovery.StartDiscoveryResponse],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, api.ErrEmptyRequest.Error())
			},
		},
		{
			name: "Request with wrong provider name",
			fields: fields{
				authz:     servicetest.NewAuthorizationStrategy(true),
				tickers:   make(map[discovery.Discoverer]*discoveryTicker),
				providers: []string{"falseProvider"},
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil[*discovery.StartDiscoveryResponse],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "provider falseProvider not known")
			},
		},
		{
			name: "Wrong permission",
			fields: fields{
				authz:     servicetest.NewAuthorizationStrategy(false, testdata.MockCloudServiceID2),
				tickers:   make(map[discovery.Discoverer]*discoveryTicker),
				providers: []string{},
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{},
			},
			want: assert.Nil[*discovery.StartDiscoveryResponse],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "access denied")
			},
		},
		{
			name: "discovery interval error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
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
			want: assert.Nil[*discovery.StartDiscoveryResponse],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "interval must be greater than zero")
			},
		},
		{
			name: "K8S authorizer error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
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
			want: assert.Nil[*discovery.StartDiscoveryResponse],
			wantErr: func(t *testing.T, gotErr error) bool {
				return assert.ErrorContains(t, gotErr, "could not authenticate to Kubernetes")
			},
		},
		{
			name: "Happy path: no discovery interval error",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
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
			want: func(t *testing.T, got *discovery.StartDiscoveryResponse) bool {
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: Azure authorizer from ENV",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
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
			want: func(t *testing.T, got *discovery.StartDiscoveryResponse) bool {
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: Azure with resource group",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
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
			want: func(t *testing.T, got *discovery.StartDiscoveryResponse) bool {
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, got)
			},
			wantErr: assert.Nil[error],
		},
		{
			name: "Happy path: CSAF with domain",
			fields: fields{
				authz:             servicetest.NewAuthorizationStrategy(true),
				tickers:           make(map[discovery.Discoverer]*discoveryTicker),
				providers:         []string{ProviderCSAF},
				discoveryInterval: time.Duration(5 * time.Minute),
			},
			args: args{
				ctx: context.Background(),
				req: &discovery.StartDiscoveryRequest{
					CsafDomain: util.Ref("clouditor.io"),
				},
			},
			want: func(t *testing.T, got *discovery.StartDiscoveryResponse) bool {
				return assert.Equal(t, &discovery.StartDiscoveryResponse{Successful: true}, got)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{
				assessmentStreams: tt.fields.assessmentStreams,
				assessment:        tt.fields.assessment,
				storage:           tt.fields.storage,
				tickers:           tt.fields.tickers,
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

			gotRes, err := svc.Start(tt.args.ctx, tt.args.req)

			tt.want(t, gotRes)
			tt.wantErr(t, err)
		})
	}
}

func TestDefaultServiceSpec(t *testing.T) {
	tests := []struct {
		name      string
		prepViper func()
		want      assert.Want[launcher.ServiceSpec]
	}{
		{
			name: "Happy path: providers given",
			prepViper: func() {
				viper.Set(config.DiscoveryProviderFlag, "azure")

			},
			want: func(t *testing.T, got launcher.ServiceSpec) bool {
				return assert.NotNil(t, got)
			},
		},
		{
			name:      "Happy path: no providers given",
			prepViper: func() {},
			want: func(t *testing.T, got launcher.ServiceSpec) bool {
				return assert.NotNil(t, got)

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			tt.prepViper()

			got := DefaultServiceSpec()

			tt.want(t, got)
		})
	}
}
