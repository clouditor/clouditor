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
		all                    bool
		CertificationTargetIDs []string
	}
	type args struct {
		ctx context.Context
		in1 service.RequestType
		req api.CertificationTargetRequest
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
				all:                    false,
				CertificationTargetIDs: []string{testdata.MockCertificationTargetID1},
			},
			args: args{
				req: &orchestrator.GetCertificationTargetRequest{CertificationTargetId: testdata.MockCertificationTargetID2},
			},
			want: false,
		},
		{
			name: "only service 1 - allowed",
			fields: fields{
				all:                    false,
				CertificationTargetIDs: []string{testdata.MockCertificationTargetID1},
			},
			args: args{
				req: &orchestrator.GetCertificationTargetRequest{CertificationTargetId: testdata.MockCertificationTargetID1},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthorizationStrategyMock{
				all:                    tt.fields.all,
				CertificationTargetIDs: tt.fields.CertificationTargetIDs,
			}
			if got := a.CheckAccess(tt.args.ctx, tt.args.in1, tt.args.req); got != tt.want {
				t.Errorf("AuthorizationStrategyMock.CheckAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
