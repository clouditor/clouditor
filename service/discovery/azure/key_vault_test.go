// Copyright 2023 Fraunhofer AISEC
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
	"fmt"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyVaultDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "No input",
			want: &azureKeyVaultDiscovery{
				azureDiscovery: &azureDiscovery{
					discovererComponent: KeyVaultComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockComputeSender{})},
			},
			want: &azureKeyVaultDiscovery{
				azureDiscovery: &azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockComputeSender{},
						},
					},
					discovererComponent: KeyVaultComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureKeyVaultDiscovery{
				azureDiscovery: &azureDiscovery{
					cred:                &mockAuthorizer{},
					discovererComponent: KeyVaultComponent,
					csID:                discovery.DefaultCloudServiceID,
					backupMap:           make(map[string]*backup),
				},
			},
		},
		{
			name: "With cloud service ID",
			args: args{
				opts: []DiscoveryOption{WithCloudServiceID(testdata.MockCloudServiceID1)},
			},
			want: &azureKeyVaultDiscovery{
				azureDiscovery: &azureDiscovery{
					discovererComponent: KeyVaultComponent,
					csID:                testdata.MockCloudServiceID1,
					backupMap:           make(map[string]*backup),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewKeyVaultDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, d)
			assert.Equal(t, "Azure Key Vault", d.Name())
			assert.Equal(t, tt.want.CloudServiceID(), d.CloudServiceID())
		})
	}
}

//type mockKeyVaultSender struct {
//	mockSender
//}
//
//func (s *mockKeyVaultSender) Do(req *http.Request) (res *http.Response, err error) {
//	if req.URL.Path == "GET https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resources?$filter=resourceType eq 'Microsoft.KeyVault/vaults'&api-version=2015-11-01" {
//		return createResponse(req, map[string]interface{}{
//			"value": &[]map[string]interface{}{
//				{
//					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.KeyVault/vaults/keyvault1",
//					"name":     "keyvault1",
//					"location": "eastus",
//				},
//				{
//					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.KeyVault/vaults/keyvault2",
//					"name":     "keyvault2",
//					"location": "westeurope",
//				},
//			},
//		}, 200)
//
//	} else {
//		// If req doesn't match, call method of anonymous field, i.e. returns error message in most cases
//		return s.mockSender.Do(req)
//	}
//}

func Test_azureKeyVaultDiscovery_List(t *testing.T) {
	// Todo 1(lebogg): Write simple test
	//d := NewKeyVaultDiscovery(WithSender(mockKeyVaultSender{}))
	//req := "GET https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resources?$filter=resourceType eq 'Microsoft.KeyVault/vaults'&api-version=2015-11-01"

	// TODO 2(lebogg): Use table
}

func Test_getKeyIDs(t *testing.T) {
	type args struct {
		keys []*voc.Key
	}
	tests := []struct {
		name       string
		args       args
		wantKeyIDs []voc.ResourceID
	}{
		{
			name: "happy path - 2 keys",
			args: args{
				keys: []*voc.Key{
					{
						Resource: &voc.Resource{ID: "key1"},
					},
					{
						Resource: &voc.Resource{ID: "key2"},
					},
				},
			},
			wantKeyIDs: []voc.ResourceID{"key1", "key2"},
		},
		{
			name: "slice of keys is empty - return empty slice of resource ids",
			args: args{
				keys: []*voc.Key{},
			},
			wantKeyIDs: []voc.ResourceID{},
		},
		{
			name: "slice of keys is nil - return empty slice of resource ids",
			args: args{
				keys: nil,
			},
			wantKeyIDs: []voc.ResourceID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantKeyIDs, getKeyIDs(tt.args.keys), "getKeyIDs(%v)", tt.args.keys)
		})
	}
}

func Test_getSecretIDs(t *testing.T) {
	type args struct {
		secrets []*voc.Secret
	}
	tests := []struct {
		name          string
		args          args
		wantSecretIDs []voc.ResourceID
	}{
		{
			name: "happy path - 2 secrets",
			args: args{
				secrets: []*voc.Secret{
					{
						Resource: &voc.Resource{ID: "Secret1"},
					},
					{
						Resource: &voc.Resource{ID: "Secret2"},
					},
				},
			},
			wantSecretIDs: []voc.ResourceID{"Secret1", "Secret2"},
		},
		{
			name: "slice of secrets is empty - return empty slice of resource ids",
			args: args{
				secrets: []*voc.Secret{},
			},
			wantSecretIDs: []voc.ResourceID{},
		},
		{
			name: "slice of secrets is nil - return empty slice of resource ids",
			args: args{
				secrets: nil,
			},
			wantSecretIDs: []voc.ResourceID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantSecretIDs, getSecretIDs(tt.args.secrets), "getSecretIDs(%v)", tt.args.secrets)
		})
	}
}

func Test_getCertificateIDs(t *testing.T) {
	type args struct {
		certificates []*voc.Certificate
	}
	tests := []struct {
		name                string
		args                args
		wantCertificatesIDs []voc.ResourceID
	}{
		{
			name: "happy path - 2 certificates",
			args: args{
				certificates: []*voc.Certificate{
					{
						Resource: &voc.Resource{ID: "certificate1"},
					},
					{
						Resource: &voc.Resource{ID: "certificate2"},
					},
				},
			},
			wantCertificatesIDs: []voc.ResourceID{"certificate1", "certificate2"},
		},
		{
			name: "slice of certificates is empty - return empty slice of resource ids",
			args: args{
				certificates: []*voc.Certificate{},
			},
			wantCertificatesIDs: []voc.ResourceID{},
		},
		{
			name: "slice of certificates is nil - return empty slice of resource ids",
			args: args{
				certificates: nil,
			},
			wantCertificatesIDs: []voc.ResourceID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantCertificatesIDs, getCertificateIDs(tt.args.certificates), "getCertificateIDs(%v)", tt.args.certificates)
		})
	}
}
func Test_azureKeyVaultDiscovery_isActive(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
		metricsClient  *azquery.MetricsClient
	}
	type args struct {
		kv *armkeyvault.Vault
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantIsActive bool
		wantErr      assert.ErrorAssertionFunc
	}{
		//{
		//	name: "happy path",
		//	fields: fields{
		//		azureDiscovery: NewMockAzureDiscovery(&mockKeyVaultSender{}),
		//	},
		//	args:         args{kv: &armkeyvault.Vault{ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.KeyVault/vaults/myKeyVault")}},
		//	wantIsActive: true,
		//	wantErr:      assert.NoError, // TODO(lebogg): Does not work yet. Since I cannot mock it currently
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureKeyVaultDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
				metricsClient:  tt.fields.metricsClient,
			}
			gotIsActive, err := d.isActive(tt.args.kv)
			if !tt.wantErr(t, err, fmt.Sprintf("isActive(%v)", tt.args.kv)) {
				return
			}
			assert.Equalf(t, tt.wantIsActive, gotIsActive, "isActive(%v)", tt.args.kv)
		})
	}
}

func Test_getCertificateName(t *testing.T) {
	certName := "SomeCertificationName"
	certIDString := "https://SomeKeyVault.vault.azure.net/certificates/" + certName
	type args struct {
		id *azcertificates.ID
	}
	tests := []struct {
		name         string
		args         args
		wantCertName string
	}{
		{
			name:         "Happy path - get name of rightly formatted ID",
			args:         args{id: util.Ref(azcertificates.ID(certIDString))},
			wantCertName: certName,
		},
		{
			name:         "Empty string provided - return empty string as well",
			args:         args{id: util.Ref(azcertificates.ID(""))},
			wantCertName: "",
		},
		{
			name:         "Wrongly formatted ID - return ID as name",
			args:         args{id: util.Ref(azcertificates.ID("subscriptions/SomeID"))},
			wantCertName: "subscriptions/SomeID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantCertName, getCertificateName(tt.args.id), "getCertificateName(%v)", tt.args.id)
		})
	}
}

func Test_convertTime(t *testing.T) {
	someDateTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	someDateToOldForUnix := time.Date(1200, time.January, 1, 0, 0, 0, 0, time.UTC)
	someDateUnix := someDateTime.Unix()
	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want assert.ValueAssertionFunc
	}{
		{
			name: "Happy path - correct time",
			args: args{util.Ref(someDateTime)},
			want: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				gotUnix, ok := i.(int64)
				assert.True(t, ok)
				return assert.Equal(t, someDateUnix, gotUnix)
			},
		},
		{
			name: "provided time is nil - return -1",
			args: args{nil},
			want: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				gotUnix, ok := i.(int64)
				assert.True(t, ok)
				return assert.Equal(t, int64(-1), gotUnix)
			},
		},
		{
			name: "provided time is before 1970 - return some negative number",
			args: args{util.Ref(someDateToOldForUnix)},
			want: func(t assert.TestingT, i interface{}, i2 ...interface{}) bool {
				gotUnix, ok := i.(int64)
				assert.True(t, ok)
				return assert.Less(t, gotUnix, int64(0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Truef(t, tt.want(t, convertTime(tt.args.t)), "convertTime(%v)", tt.args.t)
		})
	}
}

func Test_getKeyType(t *testing.T) {
	type args struct {
		kt *armkeyvault.JSONWebKeyType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Unsupported key type",
			args: args{kt: util.Ref(armkeyvault.JSONWebKeyType("NotSupportedKeyType"))},
			want: "NotSupportedKeyType",
		},
		{
			name: "EC 1",
			args: args{kt: util.Ref(armkeyvault.JSONWebKeyTypeEC)},
			want: "EC",
		},
		{
			name: "EC 2",
			args: args{kt: util.Ref(armkeyvault.JSONWebKeyTypeECHSM)},
			want: "EC",
		},
		{
			name: "RSA 1",
			args: args{kt: util.Ref(armkeyvault.JSONWebKeyTypeRSA)},
			want: "RSA",
		},
		{
			name: "RSA 2",
			args: args{kt: util.Ref(armkeyvault.JSONWebKeyTypeRSAHSM)},
			want: "RSA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getKeyType(tt.args.kt), "getKeyType(%v)", tt.args.kt)
		})
	}
}
