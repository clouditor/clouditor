package csaf

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
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
				Enabled:         true,
				Protocol:        constants.TLS,
				ProtocolVersion: 1.3,
				CipherSuites: []*ontology.CipherSuite{
					{
						MacAlgorithm:  constants.SHA_256,
						SessionCipher: constants.AES_128_GCM,
					},
				},
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
			got := d.providerTransportEncryption(tt.args.url)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_providerValidationErrors(t *testing.T) {
	type args struct {
		messages csaf.ProviderMetadataLoadMessages
	}
	tests := []struct {
		name string
		args args
		want assert.Want[[]*ontology.Error]
	}{
		{
			name: "messages given",
			args: args{
				messages: csaf.ProviderMetadataLoadMessages{
					csaf.ProviderMetadataLoadMessage{
						Message: "message1",
					},
					csaf.ProviderMetadataLoadMessage{
						Message: "message2",
					},
				},
			},
			want: func(t *testing.T, got []*ontology.Error) bool {
				want := []*ontology.Error{
					{
						Message: "message1",
					},
					{
						Message: "message2",
					},
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "no messages given",
			args: args{},
			want: assert.Nil[[]*ontology.Error],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErrs := providerValidationErrors(tt.args.messages)

			tt.want(t, gotErrs)
		})
	}
}
