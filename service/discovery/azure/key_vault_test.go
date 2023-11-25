package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
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

type mockKeyVaultSender struct {
	mockSender
}

func (s *mockKeyVaultSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "GET https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resources?$filter=resourceType eq 'Microsoft.KeyVault/vaults'&api-version=2015-11-01" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.KeyVault/vaults/keyvault1",
					"name":     "keyvault1",
					"location": "eastus",
				},
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG1/providers/Microsoft.KeyVault/vaults/keyvault2",
					"name":     "keyvault2",
					"location": "westeurope",
				},
			},
		}, 200)

	} else {
		// If req doesn't match, call method of anonymous field, i.e. returns error message in most cases
		return s.mockSender.Do(req)
	}
}

func Test_azureKeyVaultDiscovery_List(t *testing.T) {
	// Todo 1(lebogg): Write simple test
	//d := NewKeyVaultDiscovery(WithSender(mockKeyVaultSender{}))
	//req := "GET https://management.azure.com/subscriptions/00000000-0000-0000-0000-000000000000/resources?$filter=resourceType eq 'Microsoft.KeyVault/vaults'&api-version=2015-11-01"

	// TODO 2(lebogg): Use table
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

func Test_getIDs(t *testing.T) {
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
			assert.Equalf(t, tt.wantKeyIDs, getIDs(tt.args.keys), "getIDs(%v)", tt.args.keys)
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
		{
			name: "happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(&mockKeyVaultSender{}),
			},
			args:         args{kv: &armkeyvault.Vault{ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.KeyVault/vaults/myKeyVault")}},
			wantIsActive: true,
			wantErr:      assert.NoError, // TODO(lebogg): Does not work yet. Since I cannot mock it currently
		},
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
