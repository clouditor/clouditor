// Copyright 2022 Fraunhofer AISEC
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
	"testing"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/internal/testdata"
)

func TestUpdateCloudServiceRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		CloudService *CloudService
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&CloudService{
					Id: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateCloudServiceRequest{
				CloudService: tt.fields.CloudService,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("UpdateCloudServiceRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreAssessmentResultRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Result *assessment.AssessmentResult
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Result: &assessment.AssessmentResult{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &StoreAssessmentResultRequest{
				Result: tt.fields.Result,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("StoreAssessmentResultRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTargetOfEvaluationRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		TargetOfEvaluation *TargetOfEvaluation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateTargetOfEvaluationRequest{
				TargetOfEvaluation: tt.fields.TargetOfEvaluation,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("CreateTargetOfEvaluationRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCertificateRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Certificate *Certificate
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Certificate: &Certificate{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateCertificateRequest{
				Certificate: tt.fields.Certificate,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("CreateCertificateRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCertificateRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		Certificate *Certificate
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				Certificate: &Certificate{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateCertificateRequest{
				Certificate: tt.fields.Certificate,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("UpdateCertificateRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterCloudServiceRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		CloudService *CloudService
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				CloudService: &CloudService{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &RegisterCloudServiceRequest{
				CloudService: tt.fields.CloudService,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("RegisterCloudServiceRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTargetOfEvaluationRequest_GetCloudServiceId(t *testing.T) {
	type fields struct {
		TargetOfEvaluation *TargetOfEvaluation
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				TargetOfEvaluation: &TargetOfEvaluation{
					CloudServiceId: testdata.MockCloudServiceID1,
				},
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateTargetOfEvaluationRequest{
				TargetOfEvaluation: tt.fields.TargetOfEvaluation,
			}
			if got := req.GetCloudServiceId(); got != tt.want {
				t.Errorf("UpdateTargetOfEvaluationRequest.GetCloudServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}
