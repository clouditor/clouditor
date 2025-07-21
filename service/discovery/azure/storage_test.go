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
	"fmt"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_accountName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct ID",
			args: args{
				id: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account2/blobServices/default/containers/container4",
			},
			want: "account2",
		},
		{
			name: "Empty ID",
			args: args{
				id: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, accountName(tt.args.id))
		})
	}
}

func Test_azureStorageDiscovery_discoverStorageAccounts(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},

			want:    nil,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverStorageAccounts()
			if tt.wantErr != nil {
				if !tt.wantErr(t, err) {
					return
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, 9, len(got))
			}
		})
	}
}

func Test_storageAtRestEncryption(t *testing.T) {
	keySource := armstorage.KeySourceMicrosoftStorage

	type args struct {
		account *armstorage.Account
	}
	tests := []struct {
		name    string
		args    args
		want    *ontology.AtRestEncryption
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty account",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Empty KeySource",
			args: args{
				account: &armstorage.Account{
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "keySource is empty")
			},
		},
		{
			name: "No error",
			args: args{
				account: &armstorage.Account{
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
					},
				},
			},
			want: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
					ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
						Algorithm: "AES256",
						Enabled:   true,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := storageAtRestEncryption(tt.args.account)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_handleFileStorage(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	fileShareID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/fileServices/default/shares/fileshare1"
	fileShareName := "fileShare1"
	accountRegion := "eastus"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		fileshare       *armstorage.FileShareItem
		activityLogging *ontology.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*ontology.FileStorage]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty account",
			args: args{},
			want: assert.Nil[*ontology.FileStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Empty fileShareItem",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: assert.Nil[*ontology.FileStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "fileshare is nil")
			},
		},
		{
			name: "Error getting atRestEncryption properties",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
				fileshare: &armstorage.FileShareItem{
					ID: &fileShareID,
				},
			},
			want: assert.Nil[*ontology.FileStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get file storage properties for the atRestEncryption:")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID: util.Ref(accountID),
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
				fileshare: &armstorage.FileShareItem{
					ID:   &fileShareID,
					Name: &fileShareName,
				},
				activityLogging: &ontology.ActivityLogging{
					Enabled: true,
				},
			},
			want: func(t *testing.T, got *ontology.FileStorage) bool {
				want := &ontology.FileStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare1",
					Name:         fileShareName,
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: accountRegion,
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					ResourceLogging: &ontology.ResourceLogging{
						Enabled:                  false,
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled: true,
					},
				}

				assert.NotEmpty(t, got.Raw)
				got.Raw = ""
				return assert.Equal(t, want, got)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleFileStorage(tt.args.account, tt.args.fileshare, tt.args.activityLogging, "")
			if !tt.wantErr(t, err, fmt.Sprintf("handleFileStorage(%v, %v)", tt.args.account, tt.args.fileshare)) {
				return
			}

			tt.want(t, got)
		})
	}
}

func Test_generalizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty input",
			want: "",
		},
		{
			name: "Correct input",
			args: args{
				url: "https://account1.file.core.windows.net/",
			},
			want: "https://account1.[file,blob].core.windows.net/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generalizeURL(tt.args.url))
		})
	}
}

func Test_azureStorageDiscovery_handleStorageAccount(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	keySource := armstorage.KeySourceMicrosoftStorage
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountRegion := "eastus"
	minTLS := armstorage.MinimumTLSVersionTLS12
	endpointURL := "https://account1.blob.core.windows.net"
	httpsOnly := true

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		storagesList    []ontology.IsResource
		activityLogging *ontology.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*ontology.ObjectStorageService]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Account is empty",
			want: assert.Nil[*ontology.ObjectStorageService],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						MinimumTLSVersion: &minTLS,
						CreationTime:      &creationTime,
						PrimaryEndpoints: &armstorage.Endpoints{
							Blob: &endpointURL,
						},
						EnableHTTPSTrafficOnly: &httpsOnly,
					},
					Location: &accountRegion,
				},
				activityLogging: &ontology.ActivityLogging{
					Enabled: true,
				},
			},
			want: func(t *testing.T, got *ontology.ObjectStorageService) bool {
				want := &ontology.ObjectStorageService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1",
					Name:         accountName,
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: accountRegion,
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					TransportEncryption: &ontology.TransportEncryption{
						Enforced:        true,
						Enabled:         true,
						Protocol:        constants.TLS,
						ProtocolVersion: 1.2,
					},
					HttpEndpoint: &ontology.HttpEndpoint{
						Url: "https://account1.[file,blob].core.windows.net",
						TransportEncryption: &ontology.TransportEncryption{
							Enforced:        true,
							Enabled:         true,
							Protocol:        constants.TLS,
							ProtocolVersion: 1.2,
						},
					},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled: true,
					},
				}

				assert.NotEmpty(t, got.Raw)
				got.Raw = ""
				return assert.Equal(t, want, got)

			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			az := tt.fields.azureDiscovery

			got, err := az.handleStorageAccount(tt.args.account, tt.args.storagesList, tt.args.activityLogging, "")
			if !tt.wantErr(t, err) {
				return
			}

			tt.want(t, got)
		})
	}
}

func Test_handleObjectStorage(t *testing.T) {
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountRegion := "eastus"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	keySource := armstorage.KeySourceMicrosoftStorage
	containerID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1/blobServices/default/containers/container1"
	containerName := "container1"
	immutability := false
	publicAccess := armstorage.PublicAccessNone

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		container       *armstorage.ListContainerItem
		activityLogging *ontology.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[*ontology.ObjectStorage]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Account is empty",
			want: assert.Nil[*ontology.ObjectStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyStorageAccount)
			},
		},
		{
			name: "Container is empty",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: assert.Nil[*ontology.ObjectStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "container is nil")
			},
		},
		{
			name: "Error getting atRestEncryption properties",
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
				container: &armstorage.ListContainerItem{
					ID: &containerID,
				},
			},
			want: assert.Nil[*ontology.ObjectStorage],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get object storage properties for the atRestEncryption:")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
				container: &armstorage.ListContainerItem{
					ID:   &containerID,
					Name: &containerName,
					Properties: &armstorage.ContainerProperties{
						HasImmutabilityPolicy: &immutability,
						PublicAccess:          &publicAccess,
					},
				},
				activityLogging: &ontology.ActivityLogging{
					Enabled: true,
				},
			},
			want: func(t *testing.T, got *ontology.ObjectStorage) bool {
				want := &ontology.ObjectStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container1",
					Name:         containerName,
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: accountRegion,
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					Immutability: &ontology.Immutability{Enabled: false},
					ResourceLogging: &ontology.ResourceLogging{
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
						},
					},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled: true,
					},
					PublicAccess: false,
				}

				assert.NotEmpty(t, got.Raw)
				got.Raw = ""
				return assert.Equal(t, want, got)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleObjectStorage(tt.args.account, tt.args.container, tt.args.activityLogging)
			if !tt.wantErr(t, err, fmt.Sprintf("handleObjectStorage(%v, %v)", tt.args.account, tt.args.container)) {
				return
			}

			tt.want(t, got)
		})
	}
}

func Test_azureStorageDiscovery_discoverFileStorages(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		activityLogging *ontology.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				want0 := &ontology.FileStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare1",
					Name:         "fileshare1",
					CreationTime: timestamppb.New(creationTime),
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					ResourceLogging: &ontology.ResourceLogging{
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
				}
				want1 := &ontology.FileStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/fileservices/default/shares/fileshare2",
					Name:         "fileshare2",
					CreationTime: timestamppb.New(creationTime),
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					ResourceLogging: &ontology.ResourceLogging{
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
				}

				// Check if the list contains 2 entries
				assert.Equal(t, 2, len(got))

				// Check first element
				got0 := got[0].(*ontology.FileStorage)
				assert.NotEmpty(t, got0)
				got0.Raw = ""
				assert.Equal(t, want0, got0)

				// Check second element
				got1 := got[1].(*ontology.FileStorage)
				assert.NotEmpty(t, got1)
				got1.Raw = ""
				return assert.Equal(t, want1, got1)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// initialize file share client
			_ = d.initFileStorageClient()

			got, err := d.discoverFileStorages(tt.args.account, tt.args.activityLogging, "")
			if !tt.wantErr(t, err) {
				return
			}
			tt.want(t, got)
		})
	}
}

func Test_azureStorageDiscovery_discoverObjectStorages(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	accountID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"
	accountName := "account1"
	accountRegion := "eastus"
	keySource := armstorage.KeySourceMicrosoftStorage

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account         *armstorage.Account
		activityLogging *ontology.ActivityLogging
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					cred: nil,
				},
			},
			args: args{
				account: &armstorage.Account{
					ID: &accountID,
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID:   &accountID,
					Name: &accountName,
					Properties: &armstorage.AccountProperties{
						Encryption: &armstorage.Encryption{
							KeySource: &keySource,
						},
						CreationTime: &creationTime,
					},
					Location: &accountRegion,
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				want0 := &ontology.ObjectStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container1",
					Name:         "container1",
					CreationTime: timestamppb.New(creationTime),
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					Immutability: &ontology.Immutability{Enabled: false},
					ResourceLogging: &ontology.ResourceLogging{
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
						},
					},
					PublicAccess: true,
				}
				want1 := &ontology.ObjectStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1/blobservices/default/containers/container2",
					Name:         "container2",
					CreationTime: timestamppb.New(creationTime),
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.storage/storageaccounts/account1"),
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					Immutability: &ontology.Immutability{Enabled: false},
					ResourceLogging: &ontology.ResourceLogging{
						MonitoringLogDataEnabled: false,
						SecurityAlertsEnabled:    false,
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
						},
					},
					PublicAccess: true,
				}

				// Check if the list contains 2 entries
				assert.Equal(t, 2, len(got))

				// Check first element
				got0 := got[0].(*ontology.ObjectStorage)
				assert.NotEmpty(t, got0)
				got0.Raw = ""
				assert.Equal(t, want0, got0)

				// Check second element
				got1 := got[1].(*ontology.ObjectStorage)
				assert.NotEmpty(t, got1)
				got1.Raw = ""
				return assert.Equal(t, want1, got1)

			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// initialize blob container client
			_ = d.initBlobContainerClient()

			got, err := d.discoverObjectStorages(tt.args.account, tt.args.activityLogging, "")
			if !tt.wantErr(t, err) {
				return
			}
			tt.want(t, got)
		})
	}
}

func Test_azureStorageDiscovery_handleSqlServer(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		server *armsql.Server
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "error list pager",
		// 	fields: fields{
		// 		azureDiscovery: &azureDiscovery{
		// 			clients: clients{},
		// 		},
		// 	},
		// 	args: args{
		// 		server: &armsql.Server{
		// 			Location: util.Ref("eastus"),
		// 			ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
		// 			Name:     util.Ref("SQLServer1"),
		// 		},
		// 	},
		// 	want: nil,
		// 	wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
		// 		return assert.ErrorContains(t, err, "error getting next page: ")
		// 	},
		// },
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
					Properties: &armsql.ServerProperties{
						MinimalTLSVersion: util.Ref("1.2"),
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.RelationalDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1",
					Name:         "SQLServer1",
					CreationTime: nil,
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels:   make(map[string]string),
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armsql.Server\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1\",\"location\":\"eastus\",\"name\":\"SQLServer1\",\"properties\":{\"minimalTlsVersion\":\"1.2\"}}]}",
					TransportEncryption: &ontology.TransportEncryption{
						Enabled:         true,
						Enforced:        true,
						Protocol:        constants.TLS,
						ProtocolVersion: 1.2,
					},
					AnomalyDetections: []*ontology.AnomalyDetection{
						{
							Enabled: true,
							Scope:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1",
						},
					},
				},
				&ontology.DatabaseStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1/databases/sqldatabase1",
					Name:         "SqlDatabase1",
					CreationTime: nil,
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels:   make(map[string]string),
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.sql/servers/sqlserver1"),
					Raw:      "{\"*armsql.Database\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1\",\"location\":\"eastus\",\"name\":\"SqlDatabase1\",\"properties\":{\"isInfraEncryptionEnabled\":true}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleSqlServer(tt.args.server)
			if !tt.wantErr(t, err, fmt.Sprintf("handleSqlServer(%v, %v)", tt.args.server, tt.args.server)) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_anomalyDetectionEnabled(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		server *armsql.Server
		db     *armsql.Database
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error list pager",
			fields: fields{
				azureDiscovery: &azureDiscovery{
					clients: clients{},
				},
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
				},
				db: &armsql.Database{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1"),
					Name:     util.Ref("SqlDatabase1"),
					Location: util.Ref("eastus"),
					Properties: &armsql.DatabaseProperties{
						IsInfraEncryptionEnabled: util.Ref(true),
					},
				},
			},
			want: false,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get next page: ")
			},
		},
		{
			name: "Happy path: anomaly detection disabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer2"),
					Name:     util.Ref("SQLServer2"),
				},
				db: &armsql.Database{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer2/databases/SqlDatabase1"),
					Name:     util.Ref("SqlDatabase1"),
					Location: util.Ref("eastus"),
					Properties: &armsql.DatabaseProperties{
						IsInfraEncryptionEnabled: util.Ref(false),
					},
				},
			},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: anomaly detection enabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				server: &armsql.Server{
					Location: util.Ref("eastus"),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1"),
					Name:     util.Ref("SQLServer1"),
				},
				db: &armsql.Database{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Sql/servers/SQLServer1/databases/SqlDatabase1"),
					Name:     util.Ref("SqlDatabase1"),
					Location: util.Ref("eastus"),
					Properties: &armsql.DatabaseProperties{
						IsInfraEncryptionEnabled: util.Ref(true),
					},
				},
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.anomalyDetectionEnabled(tt.args.server, tt.args.db)

			tt.wantErr(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_azureStorageDiscovery_discoverCosmosDB(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(nil),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrGettingNextPage.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1",
					Name:         "DBAccount1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationEastUs,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"DBAccount1\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
				},
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb1",
					Name: "mongoDB1",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"DBAccount1\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB1\",\"location\":\"West Europe\",\"name\":\"mongoDB1\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								Enabled:   true,
								Algorithm: "",
								KeyUrl:    "https://testvault.vault.azure.net/keys/testkey/123456",
							},
						},
					},
				},
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb2",
					Name: "mongoDB2",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationEastUs,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"DBAccount1\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB2\",\"location\":\"eastus\",\"name\":\"mongoDB2\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								Enabled:   true,
								Algorithm: "",
								KeyUrl:    "https://testvault.vault.azure.net/keys/testkey/123456",
							},
						},
					},
				},
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount2",
					Name:         "DBAccount2",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationEastUs,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount2\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"DBAccount2\",\"properties\":{},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"},\"type\":\"Microsoft.DocumentDB/databaseAccounts\"}]}",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverCosmosDB()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_handleCosmosDB(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account *armcosmos.DatabaseAccountGetResults
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Cosmos DB account kind not given",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref(testdata.MockLocationWestEurope),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1"),
					Name:     util.Ref("DBAccount1"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						PublicNetworkAccess: util.Ref(armcosmos.PublicNetworkAccessEnabled),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &creationTime,
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1",
					Name:         "DBAccount1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Cosmos DB account kind not implemented: DatabaseAccountKindGlobalDocumentDB",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref(testdata.MockLocationWestEurope),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1"),
					Name:     util.Ref("DBAccount1"),
					Kind:     util.Ref(armcosmos.DatabaseAccountKindGlobalDocumentDB),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						PublicNetworkAccess: util.Ref(armcosmos.PublicNetworkAccessEnabled),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &creationTime,
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1",
					Name:         "DBAccount1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"GlobalDocumentDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Cosmos DB account kind not implemented: DatabaseAccountKindGlobalDocumentDB",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref(testdata.MockLocationWestEurope),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1"),
					Name:     util.Ref("DBAccount1"),
					Kind:     util.Ref(armcosmos.DatabaseAccountKindParse),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						PublicNetworkAccess: util.Ref(armcosmos.PublicNetworkAccessEnabled),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &creationTime,
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1",
					Name:         "DBAccount1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"Parse\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: ManagedKeyEncryption",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref(testdata.MockLocationWestEurope),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1"),
					Name:     util.Ref("DBAccount1"),
					Kind:     util.Ref(armcosmos.DatabaseAccountKindMongoDB),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						PublicNetworkAccess: util.Ref(armcosmos.PublicNetworkAccessEnabled),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &creationTime,
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1",
					Name:         "DBAccount1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
				},
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb1",
					Name: "mongoDB1",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB1\",\"location\":\"West Europe\",\"name\":\"mongoDB1\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Enabled:   true,
								Algorithm: AES256,
							},
						},
					},
				},
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb2",
					Name: "mongoDB2",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationEastUs,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB2\",\"location\":\"eastus\",\"name\":\"mongoDB2\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Enabled:   true,
								Algorithm: AES256,
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Happy path: CustomerKeyEncryption",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					Location: util.Ref(testdata.MockLocationEastUs),
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount2"),
					Name:     util.Ref("DBAccount2"),
					Kind:     util.Ref(armcosmos.DatabaseAccountKindMongoDB),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Properties: &armcosmos.DatabaseAccountGetProperties{
						KeyVaultKeyURI: util.Ref("https://testvault.vault.azure.net/keys/testkey/123456"),
					},
					SystemData: &armcosmos.SystemData{
						CreatedAt: &creationTime,
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DocumentDatabaseService{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount2",
					Name:         "DBAccount2",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount2\",\"kind\":\"MongoDB\",\"location\":\"eastus\",\"name\":\"DBAccount2\",\"properties\":{\"keyVaultKeyUri\":\"https://testvault.vault.azure.net/keys/testkey/123456\"},\"systemData\":{\"createdAt\":\"2017-05-24T13:28:53.004540398Z\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleCosmosDB(tt.args.account)
			if !tt.wantErr(t, err, fmt.Sprintf("handleCosmosDB(%v, %v)", tt.args.account, tt.args.account)) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureStorageDiscovery_discoverMongoDBDatabases(t *testing.T) {
	type fields struct {
		azureDiscovery         *azureDiscovery
		defenderProperties     map[string]*defenderProperties
		mongoDBResourcesClient bool
	}
	type args struct {
		account   *armcosmos.DatabaseAccountGetResults
		atRestEnc *ontology.AtRestEncryption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ontology.IsResource
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery:         NewMockAzureDiscovery(newMockSender()),
				mongoDBResourcesClient: true,
			},
			args: args{
				account: &armcosmos.DatabaseAccountGetResults{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1"),
					Name:     util.Ref("DBAccount1"),
					Kind:     util.Ref(armcosmos.DatabaseAccountKindMongoDB),
					Location: util.Ref(testdata.MockLocationWestEurope),
				},
				atRestEnc: &ontology.AtRestEncryption{
					Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
						ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
							Enabled:   true,
							Algorithm: AES256,
						},
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb1",
					Name: "mongoDB1",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationWestEurope,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\"}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB1\",\"location\":\"West Europe\",\"name\":\"mongoDB1\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Enabled:   true,
								Algorithm: AES256,
							},
						},
					},
				},
				&ontology.DatabaseStorage{
					Id:   "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1/mongodbdatabases/mongodb2",
					Name: "mongoDB2",
					GeoLocation: &ontology.GeoLocation{
						Region: testdata.MockLocationEastUs,
					},
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.documentdb/databaseaccounts/dbaccount1"),
					Raw:      "{\"*armcosmos.DatabaseAccountGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1\",\"kind\":\"MongoDB\",\"location\":\"West Europe\",\"name\":\"DBAccount1\"}],\"*armcosmos.MongoDBDatabaseGetResults\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.DocumentDB/databaseAccounts/DBAccount1/mongodbDatabases/mongoDB2\",\"location\":\"eastus\",\"name\":\"mongoDB2\",\"properties\":{},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Enabled:   true,
								Algorithm: AES256,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			d.defenderProperties = tt.fields.defenderProperties

			// initialize Mongo DB resources client
			if tt.fields.mongoDBResourcesClient {
				_ = d.initMongoDResourcesBClient()
			}

			got := d.discoverMongoDBDatabases(tt.args.account, tt.args.atRestEnc)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_checkTlsVersion(t *testing.T) {
	type args struct {
		version *string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "TLS version not implemented",
			args: args{
				version: util.Ref("TLS version 1.0"),
			},
			want: 0,
		},
		{
			name: "Happy path:TLS1_0",
			args: args{
				version: util.Ref("1.0"),
			},
			want: 1.0,
		},
		{
			name: "Happy path:TLS1_1",
			args: args{
				version: util.Ref("1.1"),
			},
			want: 1.1,
		},
		{
			name: "Happy path:TLS1_2",
			args: args{
				version: util.Ref("1.2"),
			},
			want: 1.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tlsVersion(tt.args.version); got != tt.want {
				t.Errorf("checkTlsVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureStorageDiscovery_getActivityLogging(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		account *armstorage.Account
	}
	tests := []struct {
		name                       string
		fields                     fields
		args                       args
		wantActivityLoggingAccount *ontology.ActivityLogging
		wantActivityLoggingBlob    *ontology.ActivityLogging
		wantActivityLoggingTable   *ontology.ActivityLogging
		wantActivityLoggingFile    *ontology.ActivityLogging
		wantRawAccount             string
		wantRawBlob                string
		wantRawTable               string
		wantRawFile                string
	}{
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				account: &armstorage.Account{
					ID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1"),
				},
			},
			wantActivityLoggingAccount: &ontology.ActivityLogging{
				Enabled:           true,
				LoggingServiceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1"},
			},
			wantActivityLoggingBlob:  nil,
			wantActivityLoggingTable: nil,
			wantActivityLoggingFile:  nil,
			wantRawBlob:              "",
			wantRawTable:             "",
			wantRawFile:              "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Init Diagnostic Settings Client
			_ = d.initDiagnosticsSettingsClient()

			gotActivityLoggingAccount, gotActivityLoggingBlob, gotActivityLoggingTable, gotActivityLoggingFile, gotRawAccount, gotRawBlob, gotRawTable, gotRawFile := d.getActivityLogging(tt.args.account)

			assert.Equal(t, tt.wantActivityLoggingAccount, gotActivityLoggingAccount)
			assert.Equal(t, tt.wantActivityLoggingBlob, gotActivityLoggingBlob)
			assert.Equal(t, tt.wantActivityLoggingTable, gotActivityLoggingTable)
			assert.Equal(t, tt.wantActivityLoggingFile, gotActivityLoggingFile)
			assert.NotEmpty(t, gotRawAccount)
			assert.Equal(t, tt.wantRawBlob, gotRawBlob)
			assert.Equal(t, tt.wantRawTable, gotRawTable)
			assert.Equal(t, tt.wantRawFile, gotRawFile)
		})
	}
}
