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
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/stretchr/testify/assert"
)

const (
	assessmentResultID1 = "11111111-1111-1111-1111-111111111111"
	assessmentResultID2 = "11111111-1111-1111-1111-111111111112"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

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

	assert.Nil(t, err)
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
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
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
			wantErr:  true,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
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
	type args struct {
		stream orchestrator.Orchestrator_StoreAssessmentResultsServer
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Store 2 assessment results to the map",
			args:    args{stream: &mockStreamer{counter: 0}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.StoreAssessmentResults(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResults() error = %v, wantErr %v", err, tt.wantErr)
				assert.Equal(t, 2, len(s.results))
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

func (m *mockStreamer) Recv() (*assessment.AssessmentResult, error) {

	if m.counter == 0 {
		m.counter++
		return &assessment.AssessmentResult{
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
		}, nil
	} else if m.counter == 1 {
		m.counter++
		return &assessment.AssessmentResult{
			Id:         assessmentResultID2,
			MetricId:   "assessmentResultMetricID2",
			EvidenceId: "11111111-1111-1111-1111-111111111112",
			Timestamp:  timestamppb.Now(),
			MetricConfiguration: &assessment.MetricConfiguration{
				TargetValue: toStruct(1.0),
				Operator:    "operator2",
				IsDefault:   false,
			},
			NonComplianceComments: "non_compliance_comment",
			Compliant:             true,
			ResourceId:            "resourceID",
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
