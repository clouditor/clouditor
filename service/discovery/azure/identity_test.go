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
	"net/http"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/stretchr/testify/assert"
)

type mockIdentitykSender struct {
	mockSender
}

func newMockIdentitySender() *mockIdentitykSender {
	m := &mockIdentitykSender{}
	return m
}

func (m mockIdentitykSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Network/networkInterfaces" {
		return createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkInterfaces/iface1",
					"name":     "iface1",
					"location": "eastus",
					"properties": map[string]interface{}{
						"networkSecurityGroup": map[string]interface{}{
							"id": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/networkSecurityGroups/nsg1",
						},
					},
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func TestNewAzureIdentityDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Empty input",
			args: args{
				opts: nil,
			},
			want: &azureIdentityDiscovery{
				&azureDiscovery{
					discovererComponent: IdentityComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
			},
		},
		{
			name: "With sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockIdentitykSender{})},
			},
			want: &azureIdentityDiscovery{
				&azureDiscovery{
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockIdentitykSender{},
						},
					},
					discovererComponent: IdentityComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
			},
		},
		{
			name: "With authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureIdentityDiscovery{
				&azureDiscovery{
					cred:                &mockAuthorizer{},
					discovererComponent: IdentityComponent,
					csID:                discovery.DefaultCloudServiceID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAzureIdentityDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, "Azure Identity", got.Name())
		})
	}
}

func Test_azureIdentityDiscovery_discoverIdentities(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []voc.IsCloudResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockIdentitySender()),
			},
			want:    []voc.IsCloudResource{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureIdentityDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got, err := d.discoverIdentities()
			tt.wantErr(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureIdentityDiscovery.discoverIdentities() = %v, want %v", got, tt.want)
			}
		})
	}
}
