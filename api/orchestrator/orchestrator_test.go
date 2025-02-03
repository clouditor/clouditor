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

func TestUpdateCertificationTargetRequest_GetCertificationTargetId(t *testing.T) {
	type fields struct {
		CertificationTarget *CertificationTarget
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&CertificationTarget{
					Id: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateCertificationTargetRequest{
				CertificationTarget: tt.fields.CertificationTarget,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("UpdateCertificationTargetRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreAssessmentResultRequest_GetCertificationTargetId(t *testing.T) {
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
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &StoreAssessmentResultRequest{
				Result: tt.fields.Result,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("StoreAssessmentResultRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateAuditScopeRequest_GetCertificationTargetId(t *testing.T) {
	type fields struct {
		AuditScope *AuditScope
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				&AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateAuditScopeRequest{
				AuditScope: tt.fields.AuditScope,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("CreateAuditScopeRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateCertificateRequest_GetCertificationTargetId(t *testing.T) {
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
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CreateCertificateRequest{
				Certificate: tt.fields.Certificate,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("CreateCertificateRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCertificateRequest_GetCertificationTargetId(t *testing.T) {
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
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateCertificateRequest{
				Certificate: tt.fields.Certificate,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("UpdateCertificateRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegisterCertificationTargetRequest_GetCertificationTargetId(t *testing.T) {
	type fields struct {
		CertificationTarget *CertificationTarget
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				CertificationTarget: &CertificationTarget{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &RegisterCertificationTargetRequest{
				CertificationTarget: tt.fields.CertificationTarget,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("RegisterCertificationTargetRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateAuditScopeRequest_GetCertificationTargetId(t *testing.T) {
	type fields struct {
		AuditScope *AuditScope
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				AuditScope: &AuditScope{
					CertificationTargetId: testdata.MockCertificationTargetID1,
				},
			},
			want: testdata.MockCertificationTargetID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateAuditScopeRequest{
				AuditScope: tt.fields.AuditScope,
			}
			if got := req.GetCertificationTargetId(); got != tt.want {
				t.Errorf("UpdateAuditScopeRequest.GetCertificationTargetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
