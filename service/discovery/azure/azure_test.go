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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"clouditor.io/clouditor/internal/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/stretchr/testify/assert"
)

type mockSender struct {
}

func (mockSender) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions" {
		res, err = createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
				},
			},
		}, 200)
	} else {
		res, err = createResponse(map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)
	}

	return
}

type mockAuthorizer struct{}

func (*mockAuthorizer) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var token azcore.AccessToken

	return token, nil
}

func createResponse(object map[string]interface{}, statusCode int) (res *http.Response, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	if err = enc.Encode(object); err != nil {
		return nil, fmt.Errorf("could not encode JSON object: %w", err)
	}

	body := io.NopCloser(buf)

	return &http.Response{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func TestGetResourceGroupName(t *testing.T) {
	accountId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3"
	result := resourceGroupName(accountId)

	assert.Equal(t, "res1", result)
}

func Test_labels(t *testing.T) {

	testValue1 := "testValue1"
	testValue2 := "testValue2"
	testValue3 := "testValue3"

	type args struct {
		tags map[string]*string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Empty map of tags",
			args: args{
				tags: map[string]*string{},
			},
			want: map[string]string{},
		},
		{
			name: "Tags are nil",
			args: args{},
			want: map[string]string{},
		},
		{
			name: "Valid tags",
			args: args{
				tags: map[string]*string{
					"testTag1": &testValue1,
					"testTag2": &testValue2,
					"testTag3": &testValue3,
				},
			},
			want: map[string]string{
				"testTag1": testValue1,
				"testTag2": testValue2,
				"testTag3": testValue3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, labels(tt.args.tags), "labels(%v)", tt.args.tags)
		})
	}
}

func Test_azureDiscovery_authorize(t *testing.T) {
	type fields struct {
		isAuthorized  bool
		sub           armsubscription.Subscription
		cred          azcore.TokenCredential
		clientOptions arm.ClientOptions
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Is authorized",
			fields: fields{
				isAuthorized: true,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, nil, err)
			},
		},
		{
			name: "No credentials configured",
			fields: fields{
				isAuthorized: false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrNoCredentialsConfigured)
			},
		},
		{
			name: "Error getting subscriptions",
			fields: fields{
				isAuthorized: false,
				cred:         &mockAuthorizer{},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrCouldNotGetSubscriptions.Error())
			},
		},
		{
			name: "Without errors",
			fields: fields{
				isAuthorized: false,
				cred:         &mockAuthorizer{},
				clientOptions: arm.ClientOptions{
					ClientOptions: policy.ClientOptions{
						Transport: mockSender{},
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &azureDiscovery{
				isAuthorized:  tt.fields.isAuthorized,
				sub:           tt.fields.sub,
				cred:          tt.fields.cred,
				clientOptions: tt.fields.clientOptions,
			}
			tt.wantErr(t, a.authorize())
		})
	}
}

func Test_initClient(t *testing.T) {
	var (
		subID      = "00000000-0000-0000-0000-000000000000"
		someError  = errors.New("some error")
		someClient = &armstorage.AccountsClient{}
	)

	type args struct {
		existingClient *armstorage.AccountsClient
		d              *azureDiscovery
		fun            ClientCreateFunc[armstorage.AccountsClient]
	}
	tests := []struct {
		name       string
		args       args
		wantClient assert.ValueAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "No error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockNetworkSender{},
						},
					},
				},
				fun: armstorage.NewAccountsClient,
			},
			wantClient: assert.NotEmpty,
			wantErr:    assert.NoError,
		},
		{
			name: "Some error, client does not exist",
			args: args{
				existingClient: nil,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockNetworkSender{},
						},
					},
				},
				fun: func(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*armstorage.AccountsClient, error) {
					return nil, someError
				},
			},
			wantClient: assert.Empty,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, someError)
			},
		},
		{
			name: "No error, client already exists",
			args: args{
				existingClient: someClient,
				d: &azureDiscovery{
					cred: &mockAuthorizer{},
					sub: armsubscription.Subscription{
						SubscriptionID: &subID,
					},
					clientOptions: arm.ClientOptions{
						ClientOptions: policy.ClientOptions{
							Transport: mockNetworkSender{},
						},
					},
				},
				fun: armstorage.NewAccountsClient,
			},
			wantClient: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				return assert.Same(t, i1, someClient)
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := initClient(tt.args.existingClient, tt.args.d, tt.args.fun)
			tt.wantErr(t, err)
			tt.wantClient(t, gotClient)
		})
	}
}

func NewMockAzureDiscovery(transport policy.Transporter, opts ...DiscoveryOption) *azureDiscovery {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := armsubscription.Subscription{
		SubscriptionID: &subID,
	}

	d := &azureDiscovery{
		cred: &mockAuthorizer{},
		sub:  sub,
		clientOptions: arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Transport: transport,
			},
		},
		csID: testutil.TestCloudService1,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}
