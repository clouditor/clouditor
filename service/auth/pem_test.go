// Copyright 2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

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
