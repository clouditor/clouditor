package servicetest

import (
	"context"
	"testing"

	"clouditor.io/clouditor/v2/api/orchestrator"
	"clouditor.io/clouditor/v2/internal/api"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/service"
)

func TestAuthorizationStrategyMock_CheckAccess(t *testing.T) {
	type fields struct {
		all                   bool
		TargetOfEvaluationIDs []string
	}
	type args struct {
		ctx context.Context
		in1 service.RequestType
		req api.TargetOfEvaluationRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "allow all",
			fields: fields{
				all: true,
			},
			want: true,
		},
		{
			name: "only service 1 - denied",
			fields: fields{
				all:                   false,
				TargetOfEvaluationIDs: []string{testdata.MockTargetOfEvaluationID1},
			},
			args: args{
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID2},
			},
			want: false,
		},
		{
			name: "only service 1 - allowed",
			fields: fields{
				all:                   false,
				TargetOfEvaluationIDs: []string{testdata.MockTargetOfEvaluationID1},
			},
			args: args{
				req: &orchestrator.GetTargetOfEvaluationRequest{TargetOfEvaluationId: testdata.MockTargetOfEvaluationID1},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyMock{
				all:                   tt.fields.all,
				TargetOfEvaluationIDs: tt.fields.TargetOfEvaluationIDs,
			}
			if got := a.CheckAccess(tt.args.ctx, tt.args.in1, tt.args.req); got != tt.want {
				t.Errorf("AuthorizationStrategyMock.CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
