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
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/fake"
	"time"
)

var FakeKeyServer = fake.KeysServer{}

// init initializes fake functions, e.g. for fake key client
func init() {
	// Set up fake ListPager for Keys
	FakeKeyServer.NewListPager = func(resourceGroupName string, vaultName string,
		options *armkeyvault.KeysClientListOptions) (resp azfake.PagerResponder[armkeyvault.KeysClientListResponse]) {
		resp.AddPage(200, armkeyvault.KeysClientListResponse{
			KeyListResult: armkeyvault.KeyListResult{
				Value: []*armkeyvault.Key{
					{
						ID:   util.Ref(string(mockKey1.ID)),
						Name: util.Ref(mockKey1.Name),
					},
					{
						ID:   util.Ref(string(mockKey2.ID)),
						Name: util.Ref(mockKey2.Name),
					},
				},
			},
		}, &azfake.AddPageOptions{})
		return
	}

	// Set up fake Get Function for Keys
	FakeKeyServer.Get = func(ctx context.Context, resourceGroupName string, vaultName string, keyName string,
		options *armkeyvault.KeysClientGetOptions) (
		resp azfake.Responder[armkeyvault.KeysClientGetResponse], errResp azfake.ErrorResponder) {
		if resourceGroupName != "MockResourceGroupName" {
			resp.SetResponse(404, armkeyvault.KeysClientGetResponse{}, &azfake.SetResponseOptions{})
			return
		}
		if keyName == mockKey1.Name {
			resp.SetResponse(200, armkeyvault.KeysClientGetResponse{
				Key: armkeyvault.Key{
					Properties: &armkeyvault.KeyProperties{
						Attributes: &armkeyvault.KeyAttributes{
							Enabled:   util.Ref(mockKey1.Enabled),
							Expires:   util.Ref(mockKey1.ExpirationDate),
							NotBefore: util.Ref(mockKey1.ActivationDate),
							Created:   util.Ref(mockKey1.ActivationDate),
						},
						KeySize: util.Ref(int32(mockKey1.KeySize)),
						Kty:     util.Ref(armkeyvault.JSONWebKeyType(mockKey1.KeyType)),
						KeyURI:  util.Ref(string(mockKey1.ID)),
					}, // Important for attributes like Dates
					ID:   util.Ref(string(mockKey1.ID)),
					Name: util.Ref(mockKey1.Name),
				},
			}, &azfake.SetResponseOptions{})
		}
		if keyName == mockKey2.Name {
			resp.SetResponse(200, armkeyvault.KeysClientGetResponse{
				Key: armkeyvault.Key{
					Properties: &armkeyvault.KeyProperties{
						Attributes: &armkeyvault.KeyAttributes{
							Enabled:   util.Ref(mockKey2.Enabled),
							Expires:   util.Ref(mockKey2.ExpirationDate),
							NotBefore: util.Ref(mockKey2.ActivationDate),
							Created:   util.Ref(mockKey2.ActivationDate),
						},
						KeySize: util.Ref(int32(mockKey2.KeySize)),
						Kty:     util.Ref(armkeyvault.JSONWebKeyType(mockKey2.KeyType)),
						KeyURI:  util.Ref(string(mockKey2.ID)),
					},
					ID:   util.Ref(string(mockKey2.ID)),
					Name: util.Ref(mockKey2.Name),
				},
			}, &azfake.SetResponseOptions{})
		}
		return
	}
}

var mockKey1 = &voc.Key{
	Resource: &voc.Resource{
		ID:           "https://mockvault.vault.azure.net/keys/mockKey1Name/",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockKey1Name",
		CreationTime: 0,
		Type:         []string{"Key"},
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(),
	ExpirationDate: time.Now().Add(24 * time.Hour).Unix(),
	KeyType:        "RSA",
	KeySize:        2048,
	NumberOfUsages: 0, // TODO(lebogg)
}

var mockKey2 = &voc.Key{
	Resource: &voc.Resource{
		ID:           "https://mockvault.vault.azure.net/keys/mockKey2Name/",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockKey2Name",
		CreationTime: time.Now().Unix(),
		Type:         []string{"Key"},
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(), // about 2 years
	ExpirationDate: time.Now().Add(24 * 30 * 24 * time.Hour).Unix(),
	KeyType:        "RSA",
	KeySize:        4096,
	NumberOfUsages: 0, // TODO(lebogg)
}
