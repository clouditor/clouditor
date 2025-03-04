// Copyright 2020-2024 Fraunhofer AISEC
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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

func Test_getAtRestEncryption(t *testing.T) {
	type args struct {
		enc *armmachinelearning.EncryptionProperty
	}
	tests := []struct {
		name          string
		args          args
		wantAtRestEnc *ontology.AtRestEncryption
	}{
		{
			name: "Empty input",
			args: args{},
			wantAtRestEnc: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
					ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
						Enabled:   true,
						Algorithm: AES256,
					},
				},
			},
		},
		{
			name: "Happy path: CustomerKeyEncryption",
			args: args{
				enc: &armmachinelearning.EncryptionProperty{
					Status: util.Ref(armmachinelearning.EncryptionStatusEnabled),
					KeyVaultProperties: &armmachinelearning.KeyVaultProperties{
						KeyVaultArmID: util.Ref("some KeyVault ID"),
					},
				},
			},
			wantAtRestEnc: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
					CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
						Enabled: true,
						KeyUrl:  "some keyvault id",
					},
				},
			},
		},
		{
			name: "Happy path: ManagedKeyEncryption",
			args: args{
				enc: &armmachinelearning.EncryptionProperty{
					Status: util.Ref(armmachinelearning.EncryptionStatusEnabled),
					KeyVaultProperties: &armmachinelearning.KeyVaultProperties{
						KeyVaultArmID: util.Ref(""),
					},
				},
			},
			wantAtRestEnc: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
					ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
						Enabled:   true,
						Algorithm: AES256,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAtRestEnc := getAtRestEncryption(tt.args.enc); !reflect.DeepEqual(gotAtRestEnc, tt.wantAtRestEnc) {
				t.Errorf("getAtRestEncryption() = %v, want %v", gotAtRestEnc, tt.wantAtRestEnc)
			}
		})
	}
}

func Test_getEncryptionStatus(t *testing.T) {
	type args struct {
		enc *armmachinelearning.EncryptionStatus
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Happy path: encryption disabled",
			args: args{
				enc: util.Ref(armmachinelearning.EncryptionStatusDisabled),
			},
			want: false,
		},
		{
			name: "Happy path: encryption enabled",
			args: args{
				enc: util.Ref(armmachinelearning.EncryptionStatusEnabled),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEncryptionStatus(tt.args.enc); got != tt.want {
				t.Errorf("getEncryptionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInternetAccessibleEndpoint(t *testing.T) {

	type args struct {
		status *armmachinelearning.PublicNetworkAccess
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "input is nil",
			args: args{},
			want: false,
		},
		{
			name: "Happy path: Enabled",
			args: args{
				status: util.Ref(armmachinelearning.PublicNetworkAccessEnabled),
			},
			want: true,
		},
		{
			name: "Happy path: Disabled",
			args: args{
				status: util.Ref(armmachinelearning.PublicNetworkAccessDisabled),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInternetAccessibleEndpoint(tt.args.status); got != tt.want {
				t.Errorf("getInternetAccessibleEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getResourceLogging(t *testing.T) {
	type args struct {
		log *string
	}
	tests := []struct {
		name string
		args args
		want *ontology.ResourceLogging
	}{
		{
			name: "Happy path: application insights disabled",
			args: args{
				log: util.Ref(""),
			},
			want: &ontology.ResourceLogging{
				Enabled: false,
			},
		},
		{
			name: "Happy path: application insights enabled",
			args: args{
				log: util.Ref("Some application insights string"),
			},
			want: &ontology.ResourceLogging{
				Enabled:           true,
				LoggingServiceIds: []string{resourceID(util.Ref("Some application insights string"))},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResourceLogging(tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResourceLogging() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getComputeStringList(t *testing.T) {
	type args struct {
		values []ontology.IsResource
	}
	tests := []struct {
		name string
		args args
		want assert.Want[[]string]
	}{
		{
			name: "Empty input",
			args: args{},
			want: func(t *testing.T, got []string) bool {
				return assert.Empty(t, got)
			},
		},
		{
			name: "Happy path",
			args: args{
				values: []ontology.IsResource{
					&ontology.VirtualMachine{
						Id: "1",
					},
					&ontology.ObjectStorage{
						Id: "2",
					},
				},
			},
			want: func(t *testing.T, got []string) bool {
				return assert.Equal(t, []string{"1", "2"}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getComputeStringList(tt.args.values)

			tt.want(t, got)
		})
	}
}
