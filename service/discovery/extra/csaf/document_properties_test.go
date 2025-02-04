package csaf

import (
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
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
		ctID   string
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
				ctID:   tt.fields.ctID,
				client: tt.fields.client,
			}
			got := d.documentChecksum(tt.args.checksumURL, tt.args.filename, tt.args.body, tt.args.algorithm, tt.args.h)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cipherSuite(t *testing.T) {
	type args struct {
		id uint16
	}
	tests := []struct {
		name string
		args args
		want *ontology.CipherSuite
	}{
		{
			name: "happy path",
			args: args{
				id: tls.TLS_AES_128_GCM_SHA256,
			},
			want: &ontology.CipherSuite{
				SessionCipher: constants.AES_128_GCM,
				MacAlgorithm:  constants.SHA_256,
			},
		},
		{
			name: "unknown id",
			args: args{
				id: 1234,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cipherSuite(tt.args.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_clientAuthenticity(t *testing.T) {
	type args struct {
		res *http.Response
	}
	tests := []struct {
		name string
		args args
		want *ontology.Authenticity
	}{
		{
			name: "happy path",
			args: args{
				res: &http.Response{
					Request: &http.Request{
						Header: http.Header{},
					},
					Header:     http.Header{},
					StatusCode: http.StatusOK,
				},
			},
			want: &ontology.Authenticity{
				Type: &ontology.Authenticity_NoAuthentication{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clientAuthenticity(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientAuthenticity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csafDiscovery_documentPGPSignature(t *testing.T) {
	type fields struct {
		domain string
		ctID   string
		client *http.Client
	}
	type args struct {
		signURL string
		body    []byte
		keyring openpgp.EntityList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantSig *ontology.DocumentSignature
	}{
		{
			name: "happy path",
			fields: fields{
				domain: goodProvider.Domain(),
				ctID:   config.DefaultCertificationTargetID,
				client: goodProvider.Client(),
			},
			args: args{
				signURL: "https://" + goodProvider.Domain() + "/.well-known/csaf/white/some-id.json.asc",
				body: func() []byte {
					res, _ := goodProvider.Client().Get("https://" + goodProvider.Domain() + "/.well-known/csaf/white/some-id")
					body, _ := io.ReadAll(res.Body)
					return body
				}(),
				keyring: goodProvider.Keyring,
			},
			wantSig: &ontology.DocumentSignature{
				Algorithm: "PGP",
				Errors:    nil,
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
			if gotSig := d.documentPGPSignature(tt.args.signURL, tt.args.body, tt.args.keyring); !reflect.DeepEqual(gotSig, tt.wantSig) {
				t.Errorf("csafDiscovery.documentPGPSignature() = %v, want %v", gotSig, tt.wantSig)
			}
		})
	}
}

func Test_fromError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name       string
		args       args
		wantErrors []*ontology.Error
	}{
		{
			name: "wrapped error",
			args: args{
				err: fmt.Errorf("this wraps multiple errors: %w, %w, %w",
					errors.New("first"),
					errors.New("second"),
					errors.New("third"),
				),
			},
			wantErrors: []*ontology.Error{
				{
					Message: "first",
				},
				{
					Message: "second",
				},
				{
					Message: "third",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrors := fromError(tt.args.err); !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("fromError() = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}
