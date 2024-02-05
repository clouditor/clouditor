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
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"

	"clouditor.io/clouditor/internal/util"

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
}

func setDiscoveryForKeyVault() *azureDiscovery {
	return &azureDiscovery{
		clientOptions: arm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: fake.NewVaultsServerTransport(&FakeVaultsServer)}},
		sub: &armsubscription.Subscription{
			ID:             util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000"),
			SubscriptionID: util.Ref("00000000-0000-0000-0000-000000000000"),
		},
	}
}
