// csaf contains a discover that discovery security advisory information from a CSAF trusted provider
package csaf

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/service/discovery/extra/csaf/providertest"
)

func Test_csafDiscovery_List(t *testing.T) {
	p, err := providertest.NewTestProvider()
	assert.NoError(t, err)
	go p.Serve()
	defer p.Close()

	type fields struct {
		domain string
		csID   string
		client *http.Client
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
				domain: "localhost:1234",
				csID:   discovery.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.ErrorContains(t, err, "could not load provider-metadata.json")
			},
		},
		{
			name: "happy path",
			fields: fields{
				domain: p.Domain(),
				client: p.Client(),
				csID:   discovery.DefaultCloudServiceID,
			},
			wantErr: func(t *testing.T, err error) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				domain: tt.fields.domain,
				client: tt.fields.client,
				csID:   tt.fields.csID,
			}
			gotList, err := d.List()
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantList, gotList)
		})
	}
}
