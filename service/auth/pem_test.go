package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"
)

func TestParseECPrivateKeyFromPEMWithPassword(t *testing.T) {
	type args struct {
		data     []byte
		password []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Private key with password",
			args: args{
				data: []byte(
					`-----BEGIN ENCRYPTED PRIVATE KEY-----
MIHsMFcGCSqGSIb3DQEFDTBKMCkGCSqGSIb3DQEFDDAcBAgTz/KWaEQ7xwICCAAw
DAYIKoZIhvcNAgkFADAdBglghkgBZQMEASoEEEoMbQeGZBq+RJGRyY2N8PwEgZAY
U36vBRn5HB8zNSic75MfpGXWRVXki1qm29G/DU+E68hksvfbJlqqpL12fQ5mbOz0
v8wNrNmehUyxEOQZlRPRdmgJJHObuOZ3Z49iWRJh26uvQLRYj0EdV9KkEKmSzxaF
1ZEAdLc369AgQGD33Ce9WGTtnROB6IIfFZULO5/wj/Ps32+T+jzZLIoGk+M/sng=
-----END ENCRYPTED PRIVATE KEY-----`),
				password: []byte("test"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseECPrivateKeyFromPEMWithPassword(tt.args.data, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptPKCS8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptPKCS8() = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func TestMarshalECPrivateKeyWithPassword(t *testing.T) {
	type args struct {
		key      *ecdsa.PrivateKey
		password []byte
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{
			name: "Marshal EC key",
			args: args{
				key:      &ecdsa.PrivateKey{D: big.NewInt(1), PublicKey: ecdsa.PublicKey{Curve: elliptic.P256(), X: big.NewInt(2), Y: big.NewInt(3)}},
				password: []byte("test"),
			},
			wantData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MarshalECPrivateKeyWithPassword(tt.args.key, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalECPrivateKeyWithPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotData, tt.wantData) {
			//	t.Errorf("MarshalECPrivateKeyWithPassword() = %v, want %v", gotData, tt.wantData)
			//}
		})
	}
}
