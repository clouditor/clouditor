// csaf contains a discover that discovery security advisory information from a CSAF trusted provider
package csaf

import (
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func Test_csafDiscovery_List(t *testing.T) {
	type fields struct {
		url  string
		csID string
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []ontology.IsResource
		wantErr  assert.WantErr
	}{
		{
			name: "fail",
			fields: fields{
				url:  "http://localhost:1234",
				csID: discovery.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not load provider-metadata.json")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				url:  tt.fields.url,
				csID: tt.fields.csID,
			}
			gotList, err := d.List()
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantList, gotList)
		})
	}
}
