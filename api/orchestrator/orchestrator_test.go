package orchestrator

import (
	"testing"
)

func TestCloudService_Validate(t *testing.T) {
	type fields struct {
		service *CloudService
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name:    "Service is nil",
			fields:  fields{service: nil},
			wantErr: ErrIsNil,
		},
		{
			name:    "Service name is empty",
			fields:  fields{service: &CloudService{}},
			wantErr: ErrNameIsMissing,
		},
		{
			name:    "Successful validation",
			fields:  fields{service: &CloudService{Name: "SomeName"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.service.Validate(); err != tt.wantErr {
				t.Errorf("Got Validate() error = %v, but want: %v", err, tt.wantErr)
			}
		})
	}
}
