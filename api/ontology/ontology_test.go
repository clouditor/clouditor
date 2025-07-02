package ontology

import (
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
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
			want: []string{"VirtualMachine", "Compute", "Infrastructure", "Resource"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResourceTypes(tt.args.r)
			assert.Equal(t, tt.want, got)
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
		{
			name: "happy path with plural",
			args: args{
				r: &Application{
					Id:         "some-id",
					Name:       "some-name",
					LibraryIds: []string{"some-library"},
					Raw:        "{}",
				},
			},
			want: []Relationship{
				{
					Property: "library",
					Value:    "some-library",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Related(tt.args.r)
			assert.Equal(t, tt.want, got)
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
		wantProps assert.Want[map[string]any]
		wantErr   assert.WantErr
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
			wantProps: func(t *testing.T, got map[string]any) bool {
				want := map[string]any{
					"activityLogging":            nil,
					"blockStorageIds":            []any{},
					"bootLogging":                nil,
					"creationTime":               "2024-01-01T00:00:00Z",
					"encryptionInUse":            nil,
					"geoLocation":                nil,
					"id":                         "my-id",
					"internetAccessibleEndpoint": false,
					"labels":                     map[string]any{},
					"name":                       "My VM",
					"description":                "",
					"networkInterfaceIds":        []any{},
					"malwareProtection":          nil,
					"osLogging":                  nil,
					"loggings":                   []any{},
					"raw":                        "",
					"redundancies":               []any{},
					"remoteAttestation":          nil,
					"resourceLogging":            nil,
					"automaticUpdates": map[string]any{
						"enabled":      false,
						"interval":     "172800s",
						"securityOnly": false,
					},
					"type":            []string{"VirtualMachine", "Compute", "Infrastructure", "Resource"},
					"usageStatistics": nil,
				}

				return assert.Equal(t, want, got)
			},
			wantErr: assert.Nil[error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProps, err := ResourceMap(tt.args.r)

			tt.wantErr(t, err)
			tt.wantProps(t, gotProps)
		})
	}
}

func TestListResourceIDs(t *testing.T) {
	type args struct {
		r []IsResource
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty input",
			args: args{},
			want: []string{},
		},
		{
			name: "Happy path",
			args: args{
				[]IsResource{
					&Account{Id: "test"},
					&Account{Id: "test2"},
				},
			},
			want: []string{"test", "test2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResourceIDs(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListResourceIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProtoResource(t *testing.T) {
	type args struct {
		resource IsResource
	}
	tests := []struct {
		name string
		args args
		want *Resource
	}{
		{
			name: "happy path",
			args: args{
				resource: &VirtualMachine{
					Id: "vm-1",
				},
			},
			want: &Resource{
				Type: &Resource_VirtualMachine{
					VirtualMachine: &VirtualMachine{
						Id: "vm-1",
					},
				},
			},
		},
		{
			name: "nil input",
			args: args{
				resource: nil,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ProtoResource(tt.args.resource))
		})
	}
}

func TestListResourceTypes(t *testing.T) {
	tests := []struct {
		name string
		want assert.Want[[]string]
	}{
		{
			name: "Happy path",
			want: func(t *testing.T, got []string) bool {
				return assert.NotEmpty(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ListResourceTypes()
			tt.want(t, got)
		})
	}
}
