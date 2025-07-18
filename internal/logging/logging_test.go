// Copyright 2023 Fraunhofer AISEC
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

package logging

import (
	"bytes"
	"fmt"
	"testing"

	"clouditor.io/clouditor/v2/api/assessment"
	"clouditor.io/clouditor/v2/api/evidence"
	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/api"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/sirupsen/logrus"
)

func TestRequestType_String(t *testing.T) {
	tests := []struct {
		name string
		r    RequestType
		want string
	}{
		{
			name: "Happy path: Assess",
			r:    Assess,
			want: "assessed",
		},
		{
			name: "Happy path: Add",
			r:    Add,
			want: "added",
		},
		{
			name: "Happy path: Create",
			r:    Create,
			want: "created",
		},
		{
			name: "Happy path: Register",
			r:    Register,
			want: "registered",
		},
		{
			name: "Happy path: Remove",
			r:    Remove,
			want: "removed",
		},
		{
			name: "Happy path: Store",
			r:    Store,
			want: "stored",
		},
		{
			name: "Happy path: Send",
			r:    Send,
			want: "sent",
		},
		{
			name: "Happy path: Update",
			r:    Update,
			want: "updated",
		},
		{
			name: "Happy path: Create",
			r:    RequestType(20),
			want: "unspecified",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("RequestType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogRequest(t *testing.T) {
	type args struct {
		level   logrus.Level
		reqType RequestType
		req     api.PayloadRequest
		params  []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Request missing",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Register,
			},
			want: "",
		},
		{
			name: "create target of evaluation",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Register,
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{},
				},
			},
			want: "level=debug msg=TargetOfEvaluation registered.\n",
		},

		{
			name: "create target of evaluation",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Register,
				req: &orchestrator.CreateTargetOfEvaluationRequest{
					TargetOfEvaluation: &orchestrator.TargetOfEvaluation{Id: testdata.MockTargetOfEvaluationID1},
				},
			},
			want: "level=debug msg=TargetOfEvaluation with ID '11111111-1111-1111-1111-111111111111' registered.\n",
		},
		{
			name: "Update AuditScope",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Update,
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						CatalogId:            testdata.MockCatalogID1,
					},
				},
			},
			want: "level=debug msg=AuditScope updated for Target of Evaluation '11111111-1111-1111-1111-111111111111'.\n",
		},
		{
			name: "Update AuditScope with params",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Update,
				req: &orchestrator.UpdateAuditScopeRequest{
					AuditScope: &orchestrator.AuditScope{
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						CatalogId:            testdata.MockCatalogID1,
					},
				},
				params: []string{fmt.Sprintf("and Catalog '%s'", testdata.MockCatalogID1)},
			},
			want: "level=debug msg=AuditScope updated for Target of Evaluation '11111111-1111-1111-1111-111111111111' and Catalog 'Catalog 1'.\n",
		},
		{
			name: "Send Evidence to queue",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Store,
				req: &assessment.AssessEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
					},
				},
				params: []string{fmt.Sprintf("back into queue for %s (%s)", "orchestrator", "localhost")},
			},
			want: "level=debug msg=Evidence with ID '11111111-1111-1111-1111-111111111111' stored for Target of Evaluation '11111111-1111-1111-1111-111111111111' back into queue for orchestrator (localhost).\n",
		},
		{
			name: "StoreEvidence()",
			args: args{
				level:   logrus.DebugLevel,
				reqType: Store,
				req: &evidence.StoreEvidenceRequest{
					Evidence: &evidence.Evidence{
						Id:                   testdata.MockEvidenceID1,
						TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1,
						ToolId:               testdata.MockEvidenceToolID1,
					},
				},
			},
			want: "level=debug msg=Evidence with ID '11111111-1111-1111-1111-111111111111' stored for Target of Evaluation '11111111-1111-1111-1111-111111111111' and Tool ID '39d85e98-c3da-11ed-afa1-0242ac120002'.\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			var log = &logrus.Entry{Logger: &logrus.Logger{Out: &buffer, Formatter: &logrus.TextFormatter{
				DisableColors:    true,
				DisableTimestamp: true,
				DisableQuote:     true,
			}, Level: logrus.DebugLevel}}
			LogRequest(log, tt.args.level, tt.args.reqType, tt.args.req, tt.args.params...)

			assert.Equal(t, tt.want, buffer.String())
		})
	}
}
