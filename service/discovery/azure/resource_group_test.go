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
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

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
		want   *ontology.Account
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				s: &armsubscription.Subscription{
					SubscriptionID: new(testdata.MockSubscriptionID),
					DisplayName:    new("Wonderful Subscription"),
					ID:             new(testdata.MockSubscriptionResourceID),
				},
			},
			want: &ontology.Account{
				Id:   new(testdata.MockSubscriptionResourceID),

				Name: new("Wonderful Subscription"),

				Raw:  new(string(`{"*armsubscription.Subscription":[{"displayName":"Wonderful Subscription","id":"/subscriptions/00000000-0000-0000-0000-000000000000","subscriptionId":"00000000-0000-0000-0000-000000000000"}]}`)),

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got := d.handleSubscription(tt.args.s)
			assert.Equal(t, tt.want, got)
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
		want   ontology.IsResource
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				rg: &armresources.ResourceGroup{
					ID:       new(testdata.MockResourceGroupID),
					Name:     new("res1"),
					Location: new("westus"),
					Tags: map[string]*string{
						"tag1Key": new("tag1"),
						"tag2Key": new("tag2"),
					},
				},
			},
			want: &ontology.ResourceGroup{
				Id:   new(testdata.MockResourceGroupID),

				Name: new("res1"),

				GeoLocation: &ontology.GeoLocation{
					Region: new("westus"),

				},
				Labels: map[string]string{
					"tag2Key": "tag2",
					"tag1Key": "tag1",
				},
				ParentId: new(testdata.MockSubscriptionResourceID),
				Raw:      new(string(`{"*armresources.ResourceGroup":[{"id":"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1","location":"westus","name":"res1","tags":{"tag1Key":"tag1","tag2Key":"tag2"}}]}`)),

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got := d.handleResourceGroup(tt.args.rg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureResourceGroupDiscovery_discoverResourceGroups(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []ontology.IsResource
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Discovery error",
			fields: fields{
				// Intentionally use wrong sender
				azureDiscovery: NewMockAzureDiscovery(nil),
			},
			wantList: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get next page: GET ")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender(),
					WithSubscription(&armsubscription.Subscription{
						DisplayName:    new("displayName"),
						ID:             new("/subscriptions/00000000-0000-0000-0000-000000000000"),
						SubscriptionID: new("00000000-0000-0000-0000-000000000000"),
					})),
			},
			wantList: []ontology.IsResource{
				&ontology.Account{
					Id:   new("/subscriptions/00000000-0000-0000-0000-000000000000"),

					Name: new("displayName"),

					Raw:  new(string(`{"*armsubscription.Subscription":[{"displayName":"displayName","id":"/subscriptions/00000000-0000-0000-0000-000000000000","subscriptionId":"00000000-0000-0000-0000-000000000000"}]}`)),

				},
				&ontology.ResourceGroup{
					Id:   new("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),

					Name: new("res1"),

					GeoLocation: &ontology.GeoLocation{
						Region: new("westus"),

					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: new("/subscriptions/00000000-0000-0000-0000-000000000000"),
					Raw:      new(string(`{"*armresources.ResourceGroup":[{"id":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1","location":"westus","name":"res1","tags":{"testKey1":"testTag1","testKey2":"testTag2"}}]}`)),

				},
				&ontology.ResourceGroup{
					Id:   new("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2"),

					Name: new("res2"),

					GeoLocation: &ontology.GeoLocation{
						Region: new("eastus"),

					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: new("/subscriptions/00000000-0000-0000-0000-000000000000"),
					Raw:      new(string(`{"*armresources.ResourceGroup":[{"id":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2","location":"eastus","name":"res2","tags":{"testKey1":"testTag1","testKey2":"testTag2"}}]}`)),

				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: with given resource group",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender(),
					WithResourceGroup("res1"),
					WithSubscription(&armsubscription.Subscription{
						DisplayName:    new("displayName"),
						ID:             new("/subscriptions/00000000-0000-0000-0000-000000000000"),
						SubscriptionID: new("00000000-0000-0000-0000-000000000000"),
					})),
			},
			wantList: []ontology.IsResource{
				&ontology.Account{
					Id:   new("/subscriptions/00000000-0000-0000-0000-000000000000"),

					Name: new("displayName"),

					Raw:  new(string(`{"*armsubscription.Subscription":[{"displayName":"displayName","id":"/subscriptions/00000000-0000-0000-0000-000000000000","subscriptionId":"00000000-0000-0000-0000-000000000000"}]}`)),

				},
				&ontology.ResourceGroup{
					Id:   new("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),

					Name: new("res1"),

					GeoLocation: &ontology.GeoLocation{
						Region: new("westus"),

					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: new("/subscriptions/00000000-0000-0000-0000-000000000000"),
					Raw:      new(string(`{"*armresources.ResourceGroup":[{"id":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1","location":"westus","name":"res1","tags":{"testKey1":"testTag1","testKey2":"testTag2"}}]}`)),

				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			gotList, err := d.discoverResourceGroups()

			assert.Equal(t, tt.wantList, gotList)
			tt.wantErr(t, err)
		})
	}
}
