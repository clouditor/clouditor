package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
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
