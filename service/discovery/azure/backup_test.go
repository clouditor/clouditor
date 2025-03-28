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
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"google.golang.org/protobuf/types/known/durationpb"
)

func Test_azureDiscovery_discoverBackupVaults(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[*azureDiscovery]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Backup vaults already discovered",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					backupMap: map[string]*backup{
						"testBackup": {
							backup:         make(map[string][]*ontology.Backup),
							backupStorages: []ontology.IsResource{},
						},
					},
				},
			},
			want:    assert.NotNil[*azureDiscovery],
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: storage account",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got *azureDiscovery) bool {
				want := []*ontology.Backup{
					{
						RetentionPeriod: durationpb.New(Duration7Days),
						Enabled:         true,
						StorageId:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222"),
						TransportEncryption: &ontology.TransportEncryption{
							Enforced:        true,
							Enabled:         true,
							ProtocolVersion: 1.2,
							Protocol:        constants.TLS,
						},
					},
				}

				return assert.Equal(t, want, got.backupMap[DataSourceTypeStorageAccountObject].backup["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"])
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: compute disk",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got *azureDiscovery) bool {
				want := []*ontology.Backup{
					{
						RetentionPeriod: durationpb.New(Duration30Days),
						Enabled:         true,
						StorageId:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222"),
						TransportEncryption: &ontology.TransportEncryption{
							Enforced:        true,
							Enabled:         true,
							ProtocolVersion: 1.2,
							Protocol:        constants.TLS,
						},
					},
				}

				return assert.Equal(t, want, got.backupMap[DataSourceTypeDisc].backup["/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"])
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			err := d.discoverBackupVaults()
			tt.wantErr(t, err)
			tt.want(t, d)
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
				azureDiscovery:       NewMockAzureDiscovery(nil),
				clientBackupInstance: true,
			},
			args: args{
				resourceGroup: "res1",
				vaultName:     "backupAccount1",
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get next page: GET")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery:       NewMockAzureDiscovery(newMockSender()),
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
				{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222"),
					Name: util.Ref("disk1-disk1-22222222-2222-2222-2222-222222222222"),
					Properties: &armdataprotection.BackupInstance{
						DataSourceInfo: &armdataprotection.Datasource{
							ResourceID:     util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"),
							DatasourceType: util.Ref("Microsoft.Compute/disks"),
						},
						PolicyInfo: &armdataprotection.PolicyInfo{
							PolicyID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupPolicies/backupPolicyDisk"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			if tt.fields.clientBackupInstance {
				// initialize backup instances client
				_ = d.initBackupInstancesClient()
			}
			got, err := d.discoverBackupInstances(tt.args.resourceGroup, tt.args.vaultName)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
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
		ctID                string
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
		wantResource ontology.IsResource
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
				ctID: testdata.MockTargetOfEvaluationID1,
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
			wantResource: &ontology.ObjectStorage{
				Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.dataprotection/backupvaults/backupaccount1/backupinstances/account1-account1-22222222-2222-2222-2222-222222222222",
				Name: "account1-account1-22222222-2222-2222-2222-222222222222",
				GeoLocation: &ontology.GeoLocation{
					Region: "westeurope",
				},
				CreationTime: nil,
				Labels:       nil,
				ParentId:     util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:          "{\"*armdataprotection.BackupInstanceResource\":[{\"properties\":{\"dataSourceInfo\":{\"datasourceType\":\"Microsoft.Storage/storageAccounts/blobServices\"}},\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/account1-account1-22222222-2222-2222-2222-222222222222\",\"name\":\"account1-account1-22222222-2222-2222-2222-222222222222\"}],\"*armdataprotection.BackupVaultResource\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1\",\"location\":\"westeurope\",\"name\":\"backupAccount1\"}]}",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: BlockStorage",
			fields: fields{
				ctID: testdata.MockTargetOfEvaluationID1,
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
			wantResource: &ontology.BlockStorage{
				Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.dataprotection/backupvaults/backupaccount1/backupinstances/disk1-disk1-22222222-2222-2222-2222-222222222222",
				Name: "disk1-disk1-22222222-2222-2222-2222-222222222222",
				GeoLocation: &ontology.GeoLocation{
					Region: "westeurope",
				},
				CreationTime: nil,
				Labels:       nil,
				ParentId:     util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:          "{\"*armdataprotection.BackupInstanceResource\":[{\"properties\":{\"dataSourceInfo\":{\"datasourceType\":\"Microsoft.Compute/disks\"}},\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1/backupInstances/disk1-disk1-22222222-2222-2222-2222-222222222222\",\"name\":\"disk1-disk1-22222222-2222-2222-2222-222222222222\"}],\"*armdataprotection.BackupVaultResource\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DataProtection/backupVaults/backupAccount1\",\"location\":\"westeurope\",\"name\":\"backupAccount1\"}]}",
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
				ctID:                tt.fields.ctID,
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
		backups []*ontology.Backup
	}
	tests := []struct {
		name string
		args args
		want []*ontology.Backup
	}{
		{
			name: "Happy path",
			args: args{
				backups: []*ontology.Backup{
					{
						Enabled:         true,
						Interval:        durationpb.New(90 * time.Hour * 24),
						RetentionPeriod: durationpb.New(100 * time.Hour * 24),
					},
				},
			},
			want: []*ontology.Backup{
				{
					Enabled:         true,
					Interval:        durationpb.New(90 * time.Hour * 24),
					RetentionPeriod: durationpb.New(100 * time.Hour * 24),
				},
			},
		},
		{
			name: "Happy path: empty input",
			args: args{},
			want: []*ontology.Backup{
				{
					Enabled:         false,
					RetentionPeriod: nil,
					Interval:        nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backupsEmptyCheck(tt.args.backups)
			assert.Equal(t, tt.want, got)
		})
	}
}
