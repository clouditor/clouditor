// Copyright 2020-2024 Fraunhofer AISEC
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
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_handleMLWorkspace(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)
	id := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/mlWorkspace"
	parent := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1"
	storage := "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/rg1/providers/microsoft.storage/storageaccounts/account1"
	applicationInsights := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.insights/components/appInsights1"
	keyVault := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Keyvault/vaults/keyVault1"

	type fields struct {
		d *azureDiscovery
	}
	type args struct {
		value       *armmachinelearning.Workspace
		computeList []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    assert.Want[ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				d: &azureDiscovery{},
			},
			args: args{
				value: &armmachinelearning.Workspace{
					Name: util.Ref("mlWorkspace"),
					ID:   util.Ref(id),
					SystemData: &armmachinelearning.SystemData{
						CreatedAt: util.Ref(creationTime),
					},
					Tags:     map[string]*string{"tag1": util.Ref("tag1"), "tag2": util.Ref("tag2")},
					Location: util.Ref("westeurope"),
					Properties: &armmachinelearning.WorkspaceProperties{
						PublicNetworkAccess: util.Ref(armmachinelearning.PublicNetworkAccessEnabled),
						ApplicationInsights: util.Ref(applicationInsights),
						Encryption: &armmachinelearning.EncryptionProperty{
							Status: util.Ref(armmachinelearning.EncryptionStatusEnabled),
							KeyVaultProperties: &armmachinelearning.KeyVaultProperties{
								KeyVaultArmID: util.Ref(keyVault),
							},
						},
						StorageAccount: util.Ref(storage),
					},
				},
			},
			want: func(t *testing.T, got ontology.IsResource) bool {
				got1 := got.(*ontology.MachineLearningService)

				want := &ontology.MachineLearningService{
					Id:                         resourceID(util.Ref(id)),
					Name:                       "mlWorkspace",
					CreationTime:               timestamppb.New(creationTime),
					GeoLocation:                &ontology.GeoLocation{Region: "westeurope"},
					Labels:                     map[string]string{"tag1": "tag1", "tag2": "tag2"},
					ParentId:                   util.Ref(parent),
					InternetAccessibleEndpoint: true,
					StorageIds:                 []string{storage},
					ComputeIds:                 []string{},
					Loggings: []*ontology.Logging{
						{
							Type: &ontology.Logging_ResourceLogging{
								ResourceLogging: &ontology.ResourceLogging{
									Enabled:           true,
									LoggingServiceIds: []string{resourceID(util.Ref(applicationInsights))},
								},
							},
						},
					},
				}

				assert.NotEmpty(t, got1.Raw)
				got1.Raw = ""

				return assert.Equal(t, want, got1)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.fields.d.handleMLWorkspace(tt.args.value, tt.args.computeList)

			tt.wantErr(t, err)
			tt.want(t, got)
		})
	}
}
