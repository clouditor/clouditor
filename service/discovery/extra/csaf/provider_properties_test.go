package csaf

import (
	"net/http"
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/csaf/providertest"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

func Test_csafDiscovery_providerTransportEncryption(t *testing.T) {
	p := providertest.NewTrustedProvider(nil,
		providertest.NewGoodIndexTxtWriter(),
		func(pmd *csaf.ProviderMetadata) {
			pmd.Publisher = &csaf.Publisher{
				Name:      util.Ref("Test Vendor"),
				Category:  util.Ref(csaf.CSAFCategoryVendor),
				Namespace: util.Ref("http://localhost"),
			}
		})
	defer p.Close()

	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ontology.TransportEncryption
	}{
		{
			name: "happy path",
			args: args{url: p.URL},
			fields: fields{
				client: p.Client(),
			},
			want: &ontology.TransportEncryption{
				Enabled: true,
			},
		},
		{
			name: "fail - bad certificate",
			args: args{url: p.URL},
			fields: fields{
				client: http.DefaultClient,
			},
			want: &ontology.TransportEncryption{
				Enabled: false,
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
			if got := d.providerTransportEncryption(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("providerTransportEncryption() = %v, want %v", got, tt.want)
			}
		})
	}
}
