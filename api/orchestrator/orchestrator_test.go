package orchestrator

import (
	"testing"
)

func TestRegisterCloudServiceRequest_Validate(t *testing.T) {
	type fields struct {
		req *RegisterCloudServiceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "Request is nil",
			fields:  fields{req: nil},
			wantErr: ErrRequestIsNil,
		},
		{
			name:    "Service is nil",
			fields:  fields{req: &RegisterCloudServiceRequest{Service: nil}},
			wantErr: ErrServiceIsNil,
		},
		{
			name:    "Service name is empty",
			fields:  fields{req: &RegisterCloudServiceRequest{Service: &CloudService{}}},
			wantErr: ErrNameIsMissing,
		},
		{
			name:    "Successful validation",
			fields:  fields{req: &RegisterCloudServiceRequest{Service: &CloudService{Name: "SomeName"}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.req.Validate(); err != tt.wantErr {
				t.Errorf("Got Validate() error = %v, but want: %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCloudServiceRequest_Validate(t *testing.T) {
	type fields struct {
		req *GetCloudServiceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "Request is nil",
			fields:  fields{req: nil},
			wantErr: ErrRequestIsNil,
		},
		{
			name:    "Service is nil",
			fields:  fields{req: &GetCloudServiceRequest{ServiceId: ""}},
			wantErr: ErrIDIsMissing,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.req.Validate(); err != tt.wantErr {
				t.Errorf("Got Validate() error = %v, but want: %v", err, tt.wantErr)
			}
		})
	}
}
