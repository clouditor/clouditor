// Copyright 2024 Fraunhofer AISEC
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

package azure

import (
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/ontology"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
)

func Test_tlsCipherSuites(t *testing.T) {
	type args struct {
		cs string
	}
	tests := []struct {
		name string
		args args
		want []*ontology.CipherSuite
	}{
		{
			name: "TLSCipherSuitesTLSAES128GCMSHA256",
			args: args{
				cs: string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
			},
			want: []*ontology.CipherSuite{
				{
					SessionCipher: "AES-128-GCM",
					MacAlgorithm:  "SHA-256",
				},
			},
		},
		{
			name: "TLSCipherSuitesTLSECDHERSAWITHAES256GCMSHA384",
			args: args{
				cs: string(armappservice.TLSCipherSuitesTLSECDHERSAWITHAES256GCMSHA384),
			},
			want: []*ontology.CipherSuite{
				{
					AuthenticationMechanism: "RSA",
					KeyExchangeAlgorithm:    "ECDHE",
					SessionCipher:           "AES-256-GCM",
					MacAlgorithm:            "SHA-384",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tlsCipherSuites(tt.args.cs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tlsCipherSuites() = %v, want %v", got, tt.want)
			}
		})
	}
}
