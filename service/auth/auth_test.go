// Copyright 2016-2020 Fraunhofer AISEC
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
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"google.golang.org/protobuf/proto"
)

func TestService_ListPublicKeys(t *testing.T) {
	type fields struct {
		apiKey *ecdsa.PrivateKey
	}
	type args struct {
		in0 context.Context
		in1 *auth.ListPublicKeysRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *auth.ListPublicResponse
		wantErr      bool
	}{
		{
			name: "List single public key",
			fields: fields{
				apiKey: &ecdsa.PrivateKey{
					PublicKey: ecdsa.PublicKey{
						Curve: elliptic.P256(),
						X:     big.NewInt(1),
						Y:     big.NewInt(2),
					},
				},
			},
			args: args{
				in0: context.TODO(),
				in1: &auth.ListPublicKeysRequest{},
			},
			wantResponse: &auth.ListPublicResponse{
				Keys: []*auth.JsonWebKey{
					{
						Kid: "1",
						Kty: "EC",
						Crv: "P-256",
						X:   "AQ",
						Y:   "Ag",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				apiKey: tt.fields.apiKey,
			}
			gotResponse, err := s.ListPublicKeys(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListPublicKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListPublicKeys() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
