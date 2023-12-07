// Copyright 2021 Fraunhofer AISEC
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
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/stretchr/testify/assert"
)

type mockResourceGroupSender struct {
	mockSender
}

func newMockResourceGroupSender() *mockResourceGroupSender {
	m := &mockResourceGroupSender{}
	return m
}

func (m mockResourceGroupSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
					"displayName":    "displayName",
				},
			},
		}, 200)
	} else if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups" {
		return createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
					"name": "res1",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"location": "westus",
				},
				{
					"id":   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2",
					"name": "res2",
					"tags": map[string]interface{}{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					"location": "eastus",
				},
			},
		}, 200)
	}

	return m.mockSender.Do(req)
}

func Test_azureResourceGroupDiscovery_handleSubscription(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		s *armsubscription.Subscription
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *voc.Account
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockResourceGroupSender()),
			},
			args: args{
				s: &armsubscription.Subscription{
					SubscriptionID: util.Ref("/subscriptions/" + discovery.DefaultCloudServiceID),
					DisplayName:    util.Ref(discovery.DefaultCloudServiceID),
				},
			},
			want: &voc.Account{
				Resource: &voc.Resource{
					ID:        voc.ResourceID("/subscriptions/" + discovery.DefaultCloudServiceID),
					ServiceID: testdata.MockCloudServiceID1,
					Name:      discovery.DefaultCloudServiceID,
					Type:      voc.AccountType,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureResourceGroupDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			if got := d.handleSubscription(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureResourceGroupDiscovery.handleSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureResourceGroupDiscovery_handleResourceGroup(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		rg *armresources.ResourceGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   voc.IsCloudResource
	}{
		{
			name: "Happy path",
			fields: fields{
				NewMockAzureDiscovery(newMockResourceGroupSender()),
			},
			args: args{
				rg: &armresources.ResourceGroup{
					ID:       util.Ref("/subscriptions/" + testdata.MockCloudServiceID1 + "/resourceGroups/res1"),
					Name:     util.Ref("res1"),
					Location: util.Ref("westus"),
					Tags: map[string]*string{
						"tag1Key": util.Ref("tag1"),
						"tag2Key": util.Ref("tag2"),
					},
				},
			},
			want: &voc.ResourceGroup{
				Resource: &voc.Resource{
					ID:        voc.ResourceID("/subscriptions/" + testdata.MockCloudServiceID1 + "/resourceGroups/res1"),
					ServiceID: testdata.MockCloudServiceID1,
					Name:      "res1",
					Type:      voc.ResourceGroupType,
					GeoLocation: voc.GeoLocation{
						Region: "westus",
					},
					Labels: map[string]string{
						"tag2Key": "tag2",
						"tag1Key": "tag1",
					},
					Parent: "/subscriptions/" + discovery.DefaultCloudServiceID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureResourceGroupDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			got := d.handleResourceGroup(tt.args.rg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureResourceGroupDiscovery_List(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []voc.IsCloudResource
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Authorize error",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCouldNotAuthenticate.Error())
			},
		},
		{
			name: "Discovery error",
			fields: fields{
				// Intentionally use wrong sender
				azureDiscovery: NewMockAzureDiscovery(newMockNetworkSender()),
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting next page: GET ")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockResourceGroupSender()),
			},
			wantList: []voc.IsCloudResource{
				&voc.Account{
					Resource: &voc.Resource{
						ID:        voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000"),
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "displayName",
						Type:      voc.AccountType,
					},
				},
				&voc.ResourceGroup{
					Resource: &voc.Resource{
						ID:        voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "res1",
						Type:      voc.ResourceGroupType,
						GeoLocation: voc.GeoLocation{
							Region: "westus",
						},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000"),
					},
				},
				&voc.ResourceGroup{
					Resource: &voc.Resource{
						ID:        voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2"),
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "res2",
						Type:      voc.ResourceGroupType,
						GeoLocation: voc.GeoLocation{
							Region: "eastus",
						},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000"),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with given resource group",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockResourceGroupSender(), WithResourceGroup("res1")),
			},
			wantList: []voc.IsCloudResource{
				&voc.Account{
					Resource: &voc.Resource{
						ID:        voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000"),
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "displayName",
						Type:      voc.AccountType,
					},
				},
				&voc.ResourceGroup{
					Resource: &voc.Resource{
						ID:        voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1"),
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "res1",
						Type:      voc.ResourceGroupType,
						GeoLocation: voc.GeoLocation{
							Region: "westus",
						},
						Labels: map[string]string{
							"testKey1": "testTag1",
							"testKey2": "testTag2",
						},
						Parent: voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000"),
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureResourceGroupDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}
			gotList, err := d.List()

			assert.Equal(t, tt.wantList, gotList)
			tt.wantErr(t, err)
		})
	}
}
