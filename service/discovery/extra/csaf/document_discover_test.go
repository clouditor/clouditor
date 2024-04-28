package csaf

import (
	"fmt"
	"net/http"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

func Test_csafDiscovery_handleAdvisory(t *testing.T) {
	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	type args struct {
		label    csaf.TLPLabel
		file     csaf.AdvisoryFile
		keyring  openpgp.EntityList
		parentId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantDoc assert.Want[*ontology.SecurityAdvisoryDocument]
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				domain: goodProvider.Domain(),
				csID:   discovery.DefaultCloudServiceID,
				client: goodProvider.Client(),
			},
			args: args{
				label: csaf.TLPLabelWhite,
				file: csaf.HashedAdvisoryFile{
					0: goodProvider.URL + "/.well-known/csaf/white/2020/some-id.json",
				},
				keyring: goodProvider.Keyring,
			},
			wantDoc: func(t *testing.T, got *ontology.SecurityAdvisoryDocument) bool {
				// Some debugging output, that can easily be used in Rego
				fmt.Println(ontology.ResourcePrettyJSON(got))
				return assert.Equal(t, "some-id", got.Id)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				domain: tt.fields.domain,
				csID:   tt.fields.csID,
				client: tt.fields.client,
			}
			gotDoc, err := d.handleAdvisory(tt.args.label, tt.args.file, tt.args.keyring, tt.args.parentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("csafDiscovery.handleAdvisory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.wantDoc(t, gotDoc)
		})
	}
}
