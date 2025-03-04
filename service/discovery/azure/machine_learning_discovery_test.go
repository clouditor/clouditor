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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
)

func Test_azureDiscovery_discoverMLWorkspaces(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error list pages",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(nil),
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSubscriptionNotFound.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, got[0].GetName(), "compute1")
				assert.Equal(t, got[1].GetName(), "mlWorkspace")
				return assert.Equal(t, 2, len(got))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverMLWorkspaces()

			tt.wantErr(t, err)
			tt.want(t, got)
		})
	}
}

func Test_azureDiscovery_discoverMLCompute(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		rg        string
		workspace *armmachinelearning.Workspace
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
				azureDiscovery: NewMockAzureDiscovery(nil),
			},
			args: args{
				rg: "rg",
				workspace: &armmachinelearning.Workspace{
					Name: util.Ref("mlWorkspace"),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, ErrSubscriptionNotFound.Error())
			},
		},
		{
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				rg: "rg1",
				workspace: &armmachinelearning.Workspace{
					Name: util.Ref("mlWorkspace"),
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				assert.Equal(t, 1, len(got))

				_, ok := got[0].(*ontology.VirtualMachine)
				return assert.True(t, ok)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverMLCompute(tt.args.rg, tt.args.workspace)

			tt.wantErr(t, err)
			tt.want(t, got)
		})
	}
}
