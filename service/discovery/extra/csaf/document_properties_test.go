package csaf

import (
	"crypto/tls"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
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
					Enabled:  true,
					Protocol: "TLS",
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
					Protocol:        "TLS",
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
					Protocol:        "TLS",
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
					Protocol:        "TLS",
				}
				return assert.Equal(t, want, got)
			},
		},
		{
			name: "state is TLS_1.3",
			args: args{
				state: &tls.ConnectionState{Version: tls.VersionTLS13},
			},
			want: func(t *testing.T, got *ontology.TransportEncryption) bool {
				want := &ontology.TransportEncryption{
					Enabled:         true,
					ProtocolVersion: 1.3,
					Protocol:        "TLS",
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
