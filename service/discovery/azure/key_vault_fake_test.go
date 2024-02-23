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

var FakeKeysServer = fake.KeysServer{}
var FakeSecretsServer = fake.SecretsServer{}

// init initializes fake functions, e.g. for fake key client
func init() {
	initFakeKeysServer()
	initFakeSecretsServer()
}

func initFakeKeysServer() {
	// Set up fake ListPager for Keys
	FakeKeysServer.NewListPager = func(resourceGroupName string, vaultName string,
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
	FakeKeysServer.Get = func(ctx context.Context, resourceGroupName string, vaultName string, keyName string,
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
					},
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
func initFakeSecretsServer() {
	// Set up fake ListPager for Secrets
	FakeSecretsServer.NewListPager = func(resourceGroupName string, vaultName string,
		options *armkeyvault.SecretsClientListOptions) (resp azfake.PagerResponder[armkeyvault.SecretsClientListResponse]) {
		resp.AddPage(200, armkeyvault.SecretsClientListResponse{
			SecretListResult: armkeyvault.SecretListResult{
				Value: []*armkeyvault.Secret{
					{
						ID:   util.Ref(string(mockSecret1.ID)),
						Name: util.Ref(mockSecret1.Name),
					},
					{
						ID:   util.Ref(string(mockSecret2.ID)),
						Name: util.Ref(mockSecret2.Name),
					},
				},
			},
		}, &azfake.AddPageOptions{})
		return
	}

	// Set up fake Get Function for Secrets
	FakeSecretsServer.Get = func(ctx context.Context, resourceGroupName string, vaultName string, keyName string,
		options *armkeyvault.SecretsClientGetOptions) (
		resp azfake.Responder[armkeyvault.SecretsClientGetResponse], errResp azfake.ErrorResponder) {
		if resourceGroupName != "MockResourceGroupName" {
			resp.SetResponse(404, armkeyvault.SecretsClientGetResponse{}, &azfake.SetResponseOptions{})
			return
		}
		if keyName == mockSecret1.Name {
			resp.SetResponse(200, armkeyvault.SecretsClientGetResponse{
				Secret: armkeyvault.Secret{
					Properties: &armkeyvault.SecretProperties{
						Attributes: &armkeyvault.SecretAttributes{
							Enabled:   util.Ref(mockKey1.Enabled),
							Expires:   util.Ref(time.Unix(mockSecret1.ExpirationDate, 0)),
							NotBefore: util.Ref(time.Unix(mockSecret1.ActivationDate, 0)),
							Created:   util.Ref(time.Unix(mockSecret1.ActivationDate, 0)),
						},
						SecretURI: util.Ref(string(mockSecret1.ID)),
					},
					ID:   util.Ref(string(mockSecret1.ID)),
					Name: util.Ref(mockSecret1.Name),
				},
			}, &azfake.SetResponseOptions{})
		}
		if keyName == mockSecret2.Name {
			resp.SetResponse(200, armkeyvault.SecretsClientGetResponse{
				Secret: armkeyvault.Secret{
					Properties: &armkeyvault.SecretProperties{
						Attributes: &armkeyvault.SecretAttributes{
							Enabled:   util.Ref(mockSecret2.Enabled),
							Expires:   util.Ref(time.Unix(mockSecret2.ExpirationDate, 0)),
							NotBefore: util.Ref(time.Unix(mockSecret2.ActivationDate, 0)),
							Created:   util.Ref(time.Unix(mockSecret2.ActivationDate, 0)),
						},
						SecretURI: util.Ref(string(mockSecret2.ID)),
					},
					ID:   util.Ref(string(mockSecret2.ID)),
					Name: util.Ref(mockSecret2.Name),
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
		Type:         voc.KeyType,
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(),
	ExpirationDate: time.Now().Add(24 * time.Hour).Unix(),
	KeyType:        "RSA",
	KeySize:        2048,
	NumberOfUsages: 0,
}

var mockKey2 = &voc.Key{
	Resource: &voc.Resource{
		ID:           "https://mockvault.vault.azure.net/keys/mockKey2Name/",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockKey2Name",
		CreationTime: time.Now().Unix(),
		Type:         voc.KeyType,
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(),
	ExpirationDate: time.Now().Add(24 * 30 * 24 * time.Hour).Unix(), // about 2 years
	KeyType:        "RSA",
	KeySize:        4096,
	NumberOfUsages: 0,
}

var mockSecret1 = &voc.Secret{
	Resource: &voc.Resource{
		ID:           "https://mockvault.vault.azure.net/secrets/mockSecret1Name/",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockSecret1Name",
		CreationTime: time.Now().Unix(),
		Type:         voc.SecretType,
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(),
	ExpirationDate: time.Now().Add(24 * time.Hour).Unix(),
	NumberOfUsages: 0,
}

var mockSecret2 = &voc.Secret{
	Resource: &voc.Resource{
		ID:           "https://mockvault.vault.azure.net/secrets/mockSecret2Name/",
		ServiceID:    "11111111-1111-1111-1111-111111111111",
		Name:         "mockSecret2Name",
		CreationTime: time.Now().Unix(),
		Type:         voc.SecretType,
		Parent:       "MockResourceGroupName",
	},
	Enabled:        true,
	ActivationDate: time.Now().Unix(),
	ExpirationDate: time.Now().Add(24 * 30 * 24 * time.Hour).Unix(), // about 2 years
	NumberOfUsages: 0,
}
