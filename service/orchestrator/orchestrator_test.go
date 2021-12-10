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

package orchestrator_test

import (
	"context"
	"io/fs"
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/persistence"

	"clouditor.io/clouditor/api/orchestrator"
	service_orchestrator "clouditor.io/clouditor/service/orchestrator"
	"github.com/stretchr/testify/assert"
)

var service *service_orchestrator.Service
var defaultTarget *orchestrator.CloudService

func TestMain(m *testing.M) {
	err := os.Chdir("../../")
	if err != nil {
		panic(err)
	}

	err = persistence.InitDB(true, "", 0)
	if err != nil {
		panic(err)
	}

	service = service_orchestrator.NewService()
	defaultTarget, err = service.CreateDefaultTargetCloudService()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestAssessmentResultHook(t *testing.T) {
	var ready1 = make(chan bool)
	var ready2 = make(chan bool)
	hookCallCounter := 0

	firstHookFunction := func(assessmentResult *orchestrator.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the firstHookFunction")

		ready1 <- true
	}

	secondHookFunction := func(assessmentResult *orchestrator.AssessmentResult, err error) {
		hookCallCounter++
		log.Println("Hello from inside the secondHookFunction")

		ready2 <- true
	}

	service := service_orchestrator.NewService()
	service.RegisterAssessmentResultHook(firstHookFunction)
	service.RegisterAssessmentResultHook(secondHookFunction)

	// Check if first hook is registered
	funcName1 := runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultsHook[0]).Pointer()).Name()
	funcName2 := runtime.FuncForPC(reflect.ValueOf(firstHookFunction).Pointer()).Name()
	assert.Equal(t, funcName1, funcName2)

	// Check if second hook is registered
	funcName1 = runtime.FuncForPC(reflect.ValueOf(service.AssessmentResultsHook[1]).Pointer()).Name()
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
					Result: &orchestrator.AssessmentResult{
						Id:                    "assessmentResultID",
						MetricId:              "assessmentResultMetricID",
						EvidenceId:            "evidenceID",
					},
				},
			},
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
		{
			name: "Store second assessment result to the map",
			args: args{
				in0: context.TODO(),
				assessment: &orchestrator.StoreAssessmentResultRequest{
					Result: &orchestrator.AssessmentResult{
						Id:                    "assessmentResultID",
						MetricId:              "assessmentResultMetricID",
						EvidenceId:            "evidenceID",
					},
				},
			},
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookCallCounter = 0
			s := service
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)
			//make the test wait
				select {
				case <- ready1:
					break
				case <-time.After(10 * time.Second):
					assert.Fail(t, "Timeout while waiting for first storeAssessmentResult to be ready")
				}

			select {
			case <- ready2:
				break
			case <-time.After(10 * time.Second):
				assert.Fail(t, "Timeout while waiting for second storeAssessmentResult to be ready")
			}


			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreAssessmentResult() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.NotEmpty(t, s.Results)
			assert.Equal(t, 2, hookCallCounter)
		})
	}


}

func TestListMetrics(t *testing.T) {
	var (
		response *orchestrator.ListMetricsResponse
		err      error
	)

	response, err = service.ListMetrics(context.TODO(), &orchestrator.ListMetricsRequest{})

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Metrics)
}

func TestListMetricConfigurations(t *testing.T) {
	var (
		response *orchestrator.ListMetricConfigurationResponse
		err      error
	)

	response, err = service.ListMetricConfigurations(context.TODO(), &orchestrator.ListMetricConfigurationRequest{})

	assert.Nil(t, err)
	assert.NotEmpty(t, response.Configurations)
}

func TestGetMetric(t *testing.T) {
	var (
		request *orchestrator.GetMetricsRequest
		metric  *assessment.Metric
		err     error
	)

	request = &orchestrator.GetMetricsRequest{
		MetricId: "TransportEncryptionEnabled",
	}

	metric, err = service.GetMetric(context.TODO(), request)

	assert.Nil(t, err)
	assert.NotNil(t, metric)
	assert.Equal(t, request.MetricId, metric.Id)
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
					Result: &orchestrator.AssessmentResult{
						Id:                    "assessmentResultID",
						MetricId:              "assessmentResultMetricID",
						EvidenceId:            "evidenceID",
					},
				},
			},
			wantErr:  false,
			wantResp: &orchestrator.StoreAssessmentResultResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service_orchestrator.NewService()
			gotResp, err := s.StoreAssessmentResult(tt.args.in0, tt.args.assessment)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreAssessmentResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("StoreAssessmentResult() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			assert.NotNil(t, s.Results["assessmentResultID"])
		})
	}
}

func TestLoad(t *testing.T) {
	var err = service_orchestrator.LoadMetrics("notfound.json")

	assert.ErrorIs(t, err, fs.ErrNotExist)

	err = service_orchestrator.LoadMetrics("metrics.json")

	assert.Nil(t, err)
}
