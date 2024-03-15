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

	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
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

func NewMockAzureDiscovery(transport policy.Transporter, opts ...DiscoveryOption) *azureDiscovery {
	var subID = "00000000-0000-0000-0000-000000000000"
	sub := &armsubscription.Subscription{
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
		csID:      testdata.MockCloudServiceID1,
		backupMap: make(map[string]*backup),
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func Test_azureDiscovery_discoverBackupVaults_Storage(t *testing.T) {
	type fields struct {
		azureDiscovery            *azureDiscovery
		clientBackupVault         bool
		clientBackupInstance      bool
		emptyClientBackupInstance bool
		clientBackupPolicy        bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "all clients missing",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockStorageSender()),
				clientBackupVault:    false,
				clientBackupInstance: false,
				clientBackupPolicy:   false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backupVaultClient and/or backupInstancesClient missing")
			},
		},
		{
			name: "backup instance client missing",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockStorageSender()),
				clientBackupVault:    true,
				clientBackupInstance: false,
				clientBackupPolicy:   false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backupVaultClient and/or backupInstancesClient missing")
			},
		},
		{
			name: "backup instance client empty",
			fields: fields{
				azureDiscovery:            NewMockAzureDiscovery(newMockStorageSender()),
				clientBackupVault:         true,
				emptyClientBackupInstance: true,
				clientBackupPolicy:        false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover backup instances:")
			},
		},

		{
			name: "Happy path",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockStorageSender()),
				clientBackupVault:    true,
				clientBackupInstance: true,
				clientBackupPolicy:   true,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				d, ok := i1.(*azureStorageDiscovery)
				if !assert.True(tt, ok) {
					return false
				}

				want := []*voc.Backup{{
					RetentionPeriod: Duration7Days,
					Enabled:         true,
					Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  constants.TLS,
					},
				},
				}

				got := d.backupMap[DataSourceTypeStorageAccountObject].backup["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"]

				return assert.Equal(t, want, got)

			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			// Set clients if needed
			if tt.fields.clientBackupVault {
				// initialize backup vaults client
				_ = d.initBackupVaultsClient()
			}
			if tt.fields.clientBackupInstance {
				// initialize backup instances client
				_ = d.initBackupInstancesClient()
			}
			if tt.fields.clientBackupPolicy {
				// initialize backup policies client
				_ = d.initBackupPoliciesClient()
			}

			// Set empty client if needed
			if tt.fields.emptyClientBackupInstance {
				// empty backup instances client
				d.clients.backupInstancesClient = &armdataprotection.BackupInstancesClient{}
			}

			err := d.discoverBackupVaults()

			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, d)
			}
		})
	}
}

func Test_azureDiscovery_discoverBackupVaults_Compute(t *testing.T) {
	type fields struct {
		azureDiscovery       *azureDiscovery
		clientBackupVault    bool
		clientBackupInstance bool
		clientBackupPolicy   bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "both clients missing",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockComputeSender()),
				clientBackupVault:    false,
				clientBackupInstance: false,
				clientBackupPolicy:   false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backupVaultClient and/or backupInstancesClient missing")
			},
		},
		{
			name: "backup instance client missing",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockComputeSender()),
				clientBackupVault:    true,
				clientBackupInstance: false,
				clientBackupPolicy:   false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "backupVaultClient and/or backupInstancesClient missing")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockComputeSender()),
				clientBackupVault:    true,
				clientBackupInstance: true,
				clientBackupPolicy:   true,
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				d, ok := i1.(*azureComputeDiscovery)
				if !assert.True(tt, ok) {
					return false
				}

				want := []*voc.Backup{{
					RetentionPeriod: Duration30Days,
					Enabled:         true,
					Storage:         voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222"),
					TransportEncryption: &voc.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  constants.TLS,
					},
				},
				}

				return assert.Equal(t, want, d.backupMap[DataSourceTypeDisc].backup["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"])
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureComputeDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			// Set clients if needed
			if tt.fields.clientBackupVault {
				// initialize backup vaults client
				_ = d.initBackupVaultsClient()
			}
			if tt.fields.clientBackupInstance {
				// initialize backup instances client
				_ = d.initBackupInstancesClient()
			}

			if tt.fields.clientBackupPolicy {
				// initialize backup policies client
				_ = d.initBackupPoliciesClient()
			}

			err := d.discoverBackupVaults()

			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, d)
			}
		})
	}
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

func Test_azureDiscovery_discoverDefender(t *testing.T) {
	type fields struct {
		azureDiscovery      *azureDiscovery
		clientDefender      bool
		emptyDefenderClient bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "defenderClient not set",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
				clientDefender: false,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "defenderClient not set")
			},
		},
		{
			name: "empty defenderClient",
			fields: fields{
				azureDiscovery:      NewMockAzureDiscovery(newMockStorageSender()),
				clientDefender:      false,
				emptyDefenderClient: true,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover pricings")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
				clientDefender: true,
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
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			// Set client if needed
			if tt.fields.clientDefender {
				// initialize backup vaults client
				_ = d.initDefenderClient()
			}

			// Set empty defender client if needed
			if tt.fields.emptyDefenderClient {
				// initialize backup vaults client
				d.clients.defenderClient = &armsecurity.PricingsClient{}
			}

			got, err := d.discoverDefender()

			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func Test_azureDiscovery_discoverBackupInstances(t *testing.T) {
	type fields struct {
		azureDiscovery       *azureDiscovery
		clientBackupInstance bool
	}
	type args struct {
		resourceGroup string
		vaultName     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*armdataprotection.BackupInstanceResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Input empty",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "missing resource group and/or vault name")
			},
		},
		{
			name: "defenderClient not set",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockNetworkSender()),
				clientBackupInstance: true,
			},
			args: args{
				resourceGroup: "res1",
				vaultName:     "backupAccount1",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting next page: GET")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockStorageSender()),
				clientBackupInstance: true,
			},
			args: args{
				resourceGroup: "res1",
				vaultName:     "backupAccount1",
			},
			wantErr: assert.NoError,
			want: []*armdataprotection.BackupInstanceResource{
				{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
					Name: util.Ref("account1-account1-22222222-2222-2222-2222-222222222222"),
					Properties: &armdataprotection.BackupInstance{
						DataSourceInfo: &armdataprotection.Datasource{
							ResourceID:     util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"),
							DatasourceType: util.Ref("Microsoft.Storage/storageAccounts/blobServices"),
						},
						PolicyInfo: &armdataprotection.PolicyInfo{
							PolicyID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyContainer"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery: tt.fields.azureDiscovery,
			}

			if tt.fields.clientBackupInstance {
				// initialize backup instances client
				_ = d.initBackupInstancesClient()
			}
			got, err := d.discoverBackupInstances(tt.args.resourceGroup, tt.args.vaultName)

			tt.wantErr(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("azureDiscovery.discoverBackupInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureDiscovery_handleInstances(t *testing.T) {
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
	}
	type args struct {
		vault    *armdataprotection.BackupVaultResource
		instance *armdataprotection.BackupInstanceResource
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResource voc.IsCloudResource
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrVaultInstanceIsEmpty.Error())
			},
		},
		{
			name: "Happy path: ObjectStorage",
			fields: fields{
				csID: testdata.MockCloudServiceID1,
			},
			args: args{
				vault: &armdataprotection.BackupVaultResource{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1"),
					Name:     util.Ref("backupAccount1"),
					Location: util.Ref("westeurope"),
				},
				instance: &armdataprotection.BackupInstanceResource{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
					Name: util.Ref("account1-account1-22222222-2222-2222-2222-222222222222"),
					Properties: &armdataprotection.BackupInstance{
						DataSourceInfo: &armdataprotection.Datasource{
							DatasourceType: util.Ref("Microsoft.Storage/storageAccounts/blobServices"),
						},
					},
				},
			},
			wantResource: &voc.ObjectStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.dataprotection/backupvaults/backupaccount1/backupinstances/account1-account1-22222222-2222-2222-2222-222222222222",
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "account1-account1-22222222-2222-2222-2222-222222222222",
						GeoLocation: voc.GeoLocation{
							Region: "westeurope",
						},
						CreationTime: 0,
						Type:         voc.ObjectStorageType,
						Labels:       nil,
						Parent:       voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
						Raw:          "{\"*armdataprotection.BackupInstanceResource\":[{\"properties\":{\"dataSourceInfo\":{\"datasourceType\":\"Microsoft.Storage/storageAccounts/blobServices\"}},\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222\",\"name\":\"account1-account1-22222222-2222-2222-2222-222222222222\"}],\"*armdataprotection.BackupVaultResource\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1\",\"location\":\"westeurope\",\"name\":\"backupAccount1\"}]}",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: BlockStorage",
			fields: fields{
				csID: testdata.MockCloudServiceID1,
			},
			args: args{
				vault: &armdataprotection.BackupVaultResource{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1"),
					Name:     util.Ref("backupAccount1"),
					Location: util.Ref("westeurope"),
				},
				instance: &armdataprotection.BackupInstanceResource{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222"),
					Name: util.Ref("disk1-disk1-22222222-2222-2222-2222-222222222222"),
					Properties: &armdataprotection.BackupInstance{
						DataSourceInfo: &armdataprotection.Datasource{
							DatasourceType: util.Ref("Microsoft.Compute/disks"),
						},
					},
				},
			},
			wantResource: &voc.BlockStorage{
				Storage: &voc.Storage{
					Resource: &voc.Resource{
						ID:        "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.dataprotection/backupvaults/backupaccount1/backupinstances/disk1-disk1-22222222-2222-2222-2222-222222222222",
						ServiceID: testdata.MockCloudServiceID1,
						Name:      "disk1-disk1-22222222-2222-2222-2222-222222222222",
						GeoLocation: voc.GeoLocation{
							Region: "westeurope",
						},
						CreationTime: 0,
						Type:         voc.BlockStorageType,
						Labels:       nil,
						Parent:       voc.ResourceID("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
						Raw:          "{\"*armdataprotection.BackupInstanceResource\":[{\"properties\":{\"dataSourceInfo\":{\"datasourceType\":\"Microsoft.Compute/disks\"}},\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222\",\"name\":\"disk1-disk1-22222222-2222-2222-2222-222222222222\"}],\"*armdataprotection.BackupVaultResource\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1\",\"location\":\"westeurope\",\"name\":\"backupAccount1\"}]}",
					},
				},
			},
			wantErr: assert.NoError,
		},
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
			}
			gotResource, err := d.handleInstances(tt.args.vault, tt.args.instance)
			tt.wantErr(t, err)

			assert.Equal(t, tt.wantResource, gotResource)
		})
	}
}

func Test_backupsEmptyCheck(t *testing.T) {
	type args struct {
		backups []*voc.Backup
	}
	tests := []struct {
		name string
		args args
		want []*voc.Backup
	}{
		{
			name: "Happy path",
			args: args{
				backups: []*voc.Backup{
					{
						Enabled:         true,
						Interval:        90,
						RetentionPeriod: 100,
					},
				},
			},
			want: []*voc.Backup{
				{
					Enabled:         true,
					Interval:        90,
					RetentionPeriod: 100,
				},
			},
		},
		{
			name: "Happy path: empty input",
			args: args{},
			want: []*voc.Backup{
				{
					Enabled:         false,
					RetentionPeriod: -1,
					Interval:        -1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := backupsEmptyCheck(tt.args.backups); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("backupsEmptyCheck() = %v, want %v", got, tt.want)
			}
		})
	}
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
			want: "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1",
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

func Test_azureStorageDiscovery_discoverDiagnosticSettings(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
	}
	type args struct {
		resourceURI string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *voc.ActivityLogging
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "No Diagnostic Setting available",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				resourceURI: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "Happy path: no workspace available",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				resourceURI: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2",
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: data logged",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockStorageSender()),
			},
			args: args{
				resourceURI: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
			},
			want: &voc.ActivityLogging{
				Logging: &voc.Logging{
					Enabled:        true,
					LoggingService: []voc.ResourceID{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1"},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &azureStorageDiscovery{
				azureDiscovery:     tt.fields.azureDiscovery,
				defenderProperties: tt.fields.defenderProperties,
			}

			// Init Diagnostic Settings Client
			_ = d.initDiagnosticsSettingsClient()

			got, err := d.discoverDiagnosticSettings(tt.args.resourceURI)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
