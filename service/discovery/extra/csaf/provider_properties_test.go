package csaf

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/gocsaf/csaf/v3/csaf"
)

func Test_csafDiscovery_providerTransportEncryption(t *testing.T) {
	type fields struct {
		domain string
		ctID   string
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
			args: args{url: goodProvider.URL},
			fields: fields{
				client: goodProvider.Client(),
			},
			want: &ontology.TransportEncryption{
				Enabled:         util.Ref(true),

				Protocol:        util.Ref(constants.TLS),

				ProtocolVersion: util.Ref(float32(1.3)),

				CipherSuites: []*ontology.CipherSuite{
					{
						MacAlgorithm:  util.Ref(constants.SHA_256),

						SessionCipher: util.Ref(constants.AES_128_GCM),

					},
				},
			},
		},
		{
			name: "fail - bad certificate",
			args: args{url: goodProvider.URL},
			fields: fields{
				client: http.DefaultClient,
			},
			want: &ontology.TransportEncryption{
				Enabled: util.Ref(false),

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				domain: tt.fields.domain,
				ctID:   tt.fields.ctID,
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
						Message: util.Ref("message1"),

					},
					{
						Message: util.Ref("message2"),

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
