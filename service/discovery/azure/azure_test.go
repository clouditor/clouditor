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
	"reflect"
	"testing"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
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
		res, err = createResponse(req, map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
				},
			},
		}, 200)
	} else {
		res, err = createResponse(req, map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)
	}

	return
}

type mockAuthorizer struct{}

func (*mockAuthorizer) GetToken(_ context.Context, _ policy.TokenRequestOptions) (azcore.AccessToken, error) {
	var token azcore.AccessToken

	return token, nil
}

func createResponse(req *http.Request, object map[string]interface{}, statusCode int) (res *http.Response, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	if err = enc.Encode(object); err != nil {
		return nil, fmt.Errorf("could not encode JSON object: %w", err)
	}

	body := io.NopCloser(buf)

	return &http.Response{
		StatusCode: statusCode,
		Body:       body,
		// We also need to fill out the request because the Azure SDK will
		// construct the error message out of this
		Request: req,
	}, nil
}

func TestNewAzureDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "Happy path",
			args: args{},
			want: &azureDiscovery{
				csID:               discovery.DefaultCloudServiceID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with sender",
			args: args{
				opts: []DiscoveryOption{WithCloudServiceID(testdata.MockCloudServiceID1)},
			},
			want: &azureDiscovery{
				csID:               testdata.MockCloudServiceID1,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with resource group",
			args: args{
				opts: []DiscoveryOption{WithResourceGroup(testdata.MockResourceGroup)},
			},
			want: &azureDiscovery{
				rg:                 util.Ref(testdata.MockResourceGroup),
				csID:               discovery.DefaultCloudServiceID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with sender",
			args: args{
				opts: []DiscoveryOption{WithSender(mockSender{})},
			},
			want: &azureDiscovery{
				clientOptions: arm.ClientOptions{
					ClientOptions: policy.ClientOptions{
						Transport: mockSender{},
					},
				},
				csID:               discovery.DefaultCloudServiceID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
		{
			name: "Happy path: with authorizer",
			args: args{
				opts: []DiscoveryOption{WithAuthorizer(&mockAuthorizer{})},
			},
			want: &azureDiscovery{
				cred:               &mockAuthorizer{},
				csID:               discovery.DefaultCloudServiceID,
				backupMap:          make(map[string]*backup),
				defenderProperties: make(map[string]*defenderProperties),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAzureDiscovery(tt.args.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TODO(anatheka): Write test
func Test_azureDiscovery_List(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		csID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
	}
	tests := []struct {
		name     string
		fields   fields
		wantList []voc.IsCloudResource
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				csID:                tt.fields.csID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			gotList, err := d.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("azureDiscovery.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotList, tt.wantList) {
				t.Errorf("azureDiscovery.List() = %v, want %v", gotList, tt.wantList)
			}
		})
	}
}

func Test_azureDiscovery_CloudServiceID(t *testing.T) {
	type fields struct {
		isAuthorized        bool
		sub                 *armsubscription.Subscription
		cred                azcore.TokenCredential
		rg                  *string
		clientOptions       arm.ClientOptions
		discovererComponent string
		clients             clients
		csID                string
		backupMap           map[string]*backup
		defenderProperties  map[string]*defenderProperties
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Happy path",
			fields: fields{
				csID: testdata.MockCloudServiceID1,
			},
			want: testdata.MockCloudServiceID1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &azureDiscovery{
				isAuthorized:        tt.fields.isAuthorized,
				sub:                 tt.fields.sub,
				cred:                tt.fields.cred,
				rg:                  tt.fields.rg,
				clientOptions:       tt.fields.clientOptions,
				discovererComponent: tt.fields.discovererComponent,
				clients:             tt.fields.clients,
				csID:                tt.fields.csID,
				backupMap:           tt.fields.backupMap,
				defenderProperties:  tt.fields.defenderProperties,
			}
			if got := a.CloudServiceID(); got != tt.want {
				t.Errorf("azureDiscovery.CloudServiceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureDiscovery_authorize(t *testing.T) {
	type fields struct {
		isAuthorized  bool
		sub           *armsubscription.Subscription
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

func TestGetResourceGroupName(t *testing.T) {
	accountId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3"
	result := resourceGroupName(accountId)

	assert.Equal(t, "res1", result)
}

func Test_resourceGroupID(t *testing.T) {
	type args struct {
		ID *string
	}
	tests := []struct {
		name string
		args args
		want voc.ResourceID
	}{
		{
			name: "invalid",
			args: args{
				ID: util.Ref("this is not a resource ID but it should not crash the Clouditor"),
			},
			want: "",
		},
		{
			name: "happy path",
			args: args{
				ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
			},
			want: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resourceGroupID(tt.args.ID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resourceGroupID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_backupPolicyName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "invalid",
			args: args{
				id: "this is not a resource ID but it should not crash the Clouditor",
			},
			want: "",
		},
		{
			name: "Happy path",
			args: args{
				id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyDisk",
			},
			want: "backupPolicyDisk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := backupPolicyName(tt.args.id); got != tt.want {
				t.Errorf("backupPolicyName() = %v, want %v", got, tt.want)
			}
		})
	}
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
					sub: &armsubscription.Subscription{
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
					sub: &armsubscription.Subscription{
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
					sub: &armsubscription.Subscription{
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

// WithDefenderProperties is a [DiscoveryOption] that adds the defender properties for our tests.
func WithDefenderProperties(dp map[string]*defenderProperties) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.defenderProperties = dp
	}
}

// WithSubscription is a [DiscoveryOption] that adds the subscription to the discoverer for our tests.
func WithSubscription(sub *armsubscription.Subscription) DiscoveryOption {
	return func(a *azureDiscovery) {
		a.sub = sub
	}
}

func NewMockAzureDiscovery(transport policy.Transporter, opts ...DiscoveryOption) *azureDiscovery {
	sub := &armsubscription.Subscription{
		SubscriptionID: util.Ref(testdata.MockSubscriptionID),
		ID:             util.Ref(testdata.MockSubscriptionResourceID),
	}

	d := &azureDiscovery{
		cred: &mockAuthorizer{},
		sub:  sub,
		clientOptions: arm.ClientOptions{
			ClientOptions: policy.ClientOptions{
				Transport: transport,
			},
		},
		csID:      testdata.MockCloudServiceID1,
		backupMap: make(map[string]*backup),
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func Test_retentionDuration(t *testing.T) {
	type args struct {
		retention string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Missing input",
			args: args{
				retention: "",
			},
			want: time.Duration(0),
		},
		{
			name: "Wrong input",
			args: args{
				retention: "TEST",
			},
			want: time.Duration(0),
		},
		{
			name: "Happy path",
			args: args{
				retention: "P30D",
			},
			want: Duration30Days,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := retentionDuration(tt.args.retention); got != tt.want {
				t.Errorf("retentionDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO(anatheka): Update test
func Test_azureDiscovery_discoverDefender(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				got, ok := i1.(map[string]*defenderProperties)
				if !assert.True(tt, ok) {
					return false
				}

				want := &defenderProperties{
					monitoringLogDataEnabled: true,
					securityAlertsEnabled:    true,
				}

				return assert.Equal(t, want, got[DefenderStorageType])
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverDefender()

			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}
