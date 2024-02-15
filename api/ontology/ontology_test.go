package ontology

import (
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/internal/util"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestResourceTypes(t *testing.T) {
	type args struct {
		r IsResource
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "happy path",
			args: args{
				r: &VirtualMachine{},
			},
			want: []string{"VirtualMachine", "Compute", "CloudResource", "Resource"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResourceTypes(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResourceTypes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelated(t *testing.T) {
	type args struct {
		r IsResource
	}
	tests := []struct {
		name string
		args args
		want []Relationship
	}{
		{
			name: "happy path",
			args: args{
				r: &ObjectStorage{
					Id:       "some-id",
					Name:     "some-name",
					ParentId: util.Ref("some-storage-account-id"),
					Raw:      "{}",
				},
			},
			want: []Relationship{
				{
					Property: "parent",
					Value:    "some-storage-account-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Related(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Related() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceMap(t *testing.T) {
	type args struct {
		r IsResource
	}
	tests := []struct {
		name      string
		args      args
		wantProps map[string]any
		wantErr   bool
	}{
		{
			name: "happy path",
			args: args{
				r: &VirtualMachine{
					Id:           "my-id",
					Name:         "My VM",
					CreationTime: timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					AutomaticUpdates: &AutomaticUpdates{
						Interval: durationpb.New(time.Hour * 24 * 2),
					},
				},
			},
			wantProps: map[string]any{
				"activityLogging":     nil,
				"blockStorageIds":     []any{},
				"bootLogging":         nil,
				"creationTime":        "2024-01-01T00:00:00Z",
				"encryptionInUse":     nil,
				"geoLocation":         nil,
				"id":                  "my-id",
				"labels":              map[string]any{},
				"name":                "My VM",
				"networkInterfaceIds": []any{},
				"malwareProtection":   nil,
				"osLogging":           nil,
				"raw":                 "",
				"resourceLogging":     nil,
				"automaticUpdates": map[string]any{
					"enabled":      false,
					"interval":     "172800s",
					"securityOnly": false,
				},
				"type": []string{"VirtualMachine", "Compute", "CloudResource", "Resource"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProps, err := ResourceMap(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantProps, gotProps)
		})
	}
}
