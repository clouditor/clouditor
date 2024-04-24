package csaf

import (
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
)

func Test_getIDsOf(t *testing.T) {
	type args struct {
		documents []ontology.IsResource
	}
	tests := []struct {
		name    string
		args    args
		wantIds []string
	}{
		{
			name:    "no documents given",
			args:    args{},
			wantIds: nil,
		},
		{
			name: "documents given",
			args: args{
				documents: []ontology.IsResource{
					&ontology.SecurityAdvisoryDocument{
						Id: "https://xx.yy.zz/XXX",
					},
				},
			},
			wantIds: []string{"https://xx.yy.zz/XXX"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIds := getIDsOf(tt.args.documents); !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("getIDsOf() = %v, want %v", gotIds, tt.wantIds)
			}
		})
	}
}
