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

package assessment

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMain(m *testing.M) {
	// make sure, that we are in the clouditor root folder to find the policies
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// TestNewService is a simply test for NewService
func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		want assessment.AssessmentServer
	}{
		{
			name: "AssessmentServer created with empty AssessmentResults map",
			want: &Service{
				AssessmentResults:             make(map[string]*assessment.Result),
				UnimplementedAssessmentServer: assessment.UnimplementedAssessmentServer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListAssessmentResults(t *testing.T) {
	var (
		response *assessment.ListAssessmentResultsResponse
		err      error
	)

	service := NewService()

	response, err = service.ListAssessmentResults(context.TODO(), &assessment.ListAssessmentResultsRequest{})

	assert.Nil(t, err)
	assert.NotNil(t, response)
}

// TestStoreEvidence tests StoreEvidence
func TestAssessEvidence(t *testing.T) {
	type args struct {
		in0      context.Context
		evidence *evidence.Evidence
	}
	tests := []struct {
		name     string
		args     args
		wantResp *evidence.StoreEvidenceResponse
		wantErr  bool
	}{
		{
			name: "Assess resource without id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStructWithTest(voc.VirtualMachine{}, t),
				},
			},
			wantErr:  true,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
		{
			name: "Assess resource without tool id",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					Timestamp: timestamppb.Now(),
					Resource:  toStructWithTest(voc.VirtualMachine{}, t),
				},
			},
			wantErr:  true,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
		{
			name: "Assess resource without timestamp",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:   "mock",
					Resource: toStructWithTest(voc.VirtualMachine{}, t),
				},
			},
			wantErr:  true,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
		{
			name: "Assess resource",
			args: args{
				in0: context.TODO(),
				evidence: &evidence.Evidence{
					ToolId:    "mock",
					Timestamp: timestamppb.Now(),
					Resource:  toStructWithTest(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}}, t),
				},
			},
			wantErr:  false,
			wantResp: &evidence.StoreEvidenceResponse{Status: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()

			_, err := s.AssessEvidence(tt.args.in0, &assessment.AssessEvidenceRequest{Evidence: tt.args.evidence})
			if (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAssessEvidences(t *testing.T) {
	type args struct {
		stream assessment.Assessment_AssessEvidencesServer
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Store 2 evidences to the map",
			args:    args{stream: &mockStreamer{counter: 0}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.AssessEvidences(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("AssessEvidences() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockStreamer struct {
	counter int
}

func (mockStreamer) SendAndClose(_ *emptypb.Empty) error {
	return nil
}

func (m *mockStreamer) Recv() (*evidence.Evidence, error) {

	resource := toStructWithoutTest(voc.VirtualMachine{Compute: &voc.Compute{CloudResource: &voc.CloudResource{ID: "my-resource-id", Type: []string{"VirtualMachine"}}}})

	if resource == nil {
		return nil, errors.New("error getting struct")
	}

	if m.counter == 0 {
		m.counter++
		return &evidence.Evidence{
			Id:        "MockEvidenceId-1",
			ServiceId: "MockServiceId-1",
			Timestamp: timestamppb.Now(),
			Raw:       "",
			Resource:  resource,
			ToolId:    "MockToolId-1",
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &evidence.Evidence{
			Id:        "MockEvidenceId-2",
			ServiceId: "MockServiceId-2",
			Timestamp: timestamppb.Now(),
			Raw:       "",
			Resource:  resource,
			ToolId:    "MockToolId-1",
		}, nil
	} else {
		return nil, io.EOF
	}
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

func toStructWithTest(r voc.IsCloudResource, t *testing.T) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		assert.NotNil(t, err)
	}

	return
}

func toStructWithoutTest(r voc.IsCloudResource) (s *structpb.Value) {
	s, err := voc.ToStruct(r)
	if err != nil {
		fmt.Println("error moving to struct: %w", err)
		return nil
	}

	return
}
