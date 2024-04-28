package csaf

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"hash"
	"net/http"
	"net/http/httptest"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
)

func Test_documentValidationErrors(t *testing.T) {
	type args struct {
		messages []string
	}
	tests := []struct {
		name string
		args args
		want assert.Want[[]*ontology.Error]
	}{
		{
			name: "messages given",
			args: args{
				messages: []string{"message1", "message2"},
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
			gotErrs := documentValidationErrors(tt.args.messages)

			tt.want(t, gotErrs)
		})
	}
}

func Test_transportEncryption(t *testing.T) {
	type args struct {
		state *tls.ConnectionState
	}
	tests := []struct {
		name string
		args args
		want assert.Want[*ontology.TransportEncryption]
	}{
		{
			name: "state is nil",
			args: args{
				state: nil,
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{Enabled: false}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state not known",
			args: args{
				state: &tls.ConnectionState{Version: 123},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:      true,
					Protocol:     constants.TLS,
					CipherSuites: []*ontology.CipherSuite{},
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state is TLS_1.0",
			args: args{
				state: &tls.ConnectionState{Version: tls.VersionTLS10},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:         true,
					ProtocolVersion: 1.0,
					Protocol:        constants.TLS,
					CipherSuites:    []*ontology.CipherSuite{},
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state is TLS_1.1",
			args: args{
				state: &tls.ConnectionState{Version: tls.VersionTLS11},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:         true,
					ProtocolVersion: 1.1,
					Protocol:        constants.TLS,
					CipherSuites:    []*ontology.CipherSuite{},
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state is TLS_1.2",
			args: args{
				state: &tls.ConnectionState{Version: tls.VersionTLS12},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:         true,
					ProtocolVersion: 1.2,
					Protocol:        constants.TLS,
					CipherSuites:    []*ontology.CipherSuite{},
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state is TLS_1.3",
			args: args{
				state: &tls.ConnectionState{
					Version:     tls.VersionTLS13,
					CipherSuite: tls.TLS_AES_256_GCM_SHA384,
				},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:         true,
					ProtocolVersion: 1.3,
					Protocol:        constants.TLS,
					CipherSuites: []*ontology.CipherSuite{
						{
							SessionCipher: "AES-256-GCM",
							MacAlgorithm:  "SHA-384",
						},
					},
				}
				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transportEncryption(tt.args.state)
			tt.want(t, got)
		})
	}
}

func Test_csafDiscovery_documentChecksum(t *testing.T) {
	var badChecksumSrv = func() *httptest.Server {
		mux := http.NewServeMux()
		mux.HandleFunc("/file.json.sha256", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("123456 file.json"))
		})
		srv := httptest.NewTLSServer(mux)
		srv.EnableHTTP2 = true
		return srv
	}()
	defer badChecksumSrv.Close()

	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	type args struct {
		checksumURL string
		filename    string
		body        []byte
		algorithm   string
		h           hash.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ontology.DocumentChecksum
	}{
		{
			name: "checksum mismatch",
			fields: fields{
				client: badChecksumSrv.Client(),
			},
			args: args{
				checksumURL: "https://" + badChecksumSrv.Listener.Addr().String() + "/file.json.sha256",
				filename:    "file.json",
				algorithm:   constants.SHA_256,
				h:           sha256.New(),
			},
			want: &ontology.DocumentChecksum{
				Errors:    fromError(errors.New("checksum mismatch")),
				Algorithm: constants.SHA_256,
			},
		},
		{
			name: "filename mismatch",
			fields: fields{
				client: badChecksumSrv.Client(),
			},
			args: args{
				checksumURL: "https://" + badChecksumSrv.Listener.Addr().String() + "/file.json.sha256",
				filename:    "anotherfile.json",
				algorithm:   constants.SHA_256,
				h:           sha256.New(),
			},
			want: &ontology.DocumentChecksum{
				Errors:    fromError(errors.New("checksum file does not contain correct filename")),
				Algorithm: constants.SHA_256,
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
			got := d.documentChecksum(tt.args.checksumURL, tt.args.filename, tt.args.body, tt.args.algorithm, tt.args.h)
			assert.Equal(t, tt.want, got)
		})
	}
}
