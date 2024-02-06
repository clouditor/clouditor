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
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"context"
	"errors"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/fake"
)

// TODO(lebogg): Add helpers for tests, e.g. Fake Server
var FakeVaultsServer = fake.VaultsServer{}

func initKeyVaultTests() {
	FakeVaultsServer.Get = func(ctx context.Context,
		resourceGroupName string,
		vaultName string,
		options *armkeyvault.VaultsClientGetOptions) (
		resp azfake.Responder[armkeyvault.VaultsClientGetResponse], errResp azfake.ErrorResponder) {
		if resourceGroupName == "fake-resource-group" {
			resp.SetResponse(200, armkeyvault.VaultsClientGetResponse{
				Vault: armkeyvault.Vault{
					Properties: nil,
					Location:   nil,
					Tags:       nil,
					ID:         nil,
					Name:       util.Ref("Fake-KeyVault-Name"),
					SystemData: nil,
					Type:       nil,
				}}, nil)
			return
		} else {
			// If no/wrong RG is set, return 404
			resp.SetResponse(404, armkeyvault.VaultsClientGetResponse{}, nil)
			return
		}
	}
	// TODO(lebogg): FakeVaultsServer.NewListPager
	FakeVaultsServer.NewListPager = func(options *armkeyvault.VaultsClientListOptions) (resp azfake.PagerResponder[armkeyvault.VaultsClientListResponse]) {
		resp.AddPage(200, armkeyvault.VaultsClientListResponse{
			ResourceListResult: armkeyvault.ResourceListResult{
				Value: []*armkeyvault.Resource{
					{
						ID:   util.Ref(string(mockKeyVault1.ID)),
						Name: util.Ref(mockKeyVault1.Name),
					},
					{
						ID:   util.Ref(string(mockKeyVault2.ID)),
						Name: util.Ref(mockKeyVault2.Name),
					},
				},
			}}, nil)
		return
	}
	// TODO(lebogg): NewListByResourceGroupPager
	FakeVaultsServer.NewListByResourceGroupPager = func(rg string, options *armkeyvault.VaultsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armkeyvault.VaultsClientListByResourceGroupResponse]) {
		if rg == string(mockKeyVault1.Parent) {
			resp.AddPage(200, armkeyvault.VaultsClientListByResourceGroupResponse{
				VaultListResult: armkeyvault.VaultListResult{
					Value: []*armkeyvault.Vault{
						{
							ID:   util.Ref(string(mockKeyVault1.ID)),
							Name: util.Ref(mockKeyVault1.Name),
						},
					},
				}}, nil)
			return
		} else { // assume wrong RG
			// TODO(lebogg): Check error. Maybe it is thrown on the "wrong" level (of pagers)
			resp.AddError(errors.New("invalid resource group"))
			return
		}
	}
}

var mockKeyVault1 = &voc.KeyVault{
	Resource: &voc.Resource{
		ID:           "",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockKeyVault1",
		CreationTime: 0,
		Type:         []string{"KeyVault", "Resource"},
		GeoLocation:  voc.GeoLocation{},
		Labels:       nil,
		Raw:          "",
		Parent:       "resource-group-1",
	},
	IsActive:     false,
	Keys:         nil,
	PublicAccess: false,
}
var mockKeyVault2 = &voc.KeyVault{
	Resource: &voc.Resource{
		ID:           "",
		ServiceID:    "",
		Name:         "mockKeyVault2",
		CreationTime: 0,
		Type:         nil,
		GeoLocation:  voc.GeoLocation{},
		Labels:       nil,
		Raw:          "",
		Parent:       "resource-group-2",
	},
	IsActive:     false,
	Keys:         nil,
	PublicAccess: false,
}
