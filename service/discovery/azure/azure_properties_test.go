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

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
)

func Test_tlsCipherSuites(t *testing.T) {
	type args struct {
		cs string
	}
	tests := []struct {
		name string
		args args
		want []*ontology.CipherSuite
	}{
		{
			name: "TLSCipherSuitesTLSAES128GCMSHA256",
			args: args{
				cs: string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
			},
			want: []*ontology.CipherSuite{
				{
					SessionCipher: "AES-128-GCM",
					MacAlgorithm:  "SHA-256",
				},
			},
		},
		{
			name: "TLSCipherSuitesTLSECDHERSAWITHAES256GCMSHA384",
			args: args{
				cs: string(armappservice.TLSCipherSuitesTLSECDHERSAWITHAES256GCMSHA384),
			},
			want: []*ontology.CipherSuite{
				{
					AuthenticationMechanism: "RSA",
					KeyExchangeAlgorithm:    "ECDHE",
					SessionCipher:           "AES-256-GCM",
					MacAlgorithm:            "SHA-384",
				},
			},
		},
		{
			name: "not a TLS cipher",
			args: args{
				cs: "NOTTLS_AES_256",
			},
			want: nil,
		},
		{
			name: "invalid authentication",
			args: args{
				cs: "TLS_ECDHE_RSB_WITH_AES_256_GCM_SHA384",
			},
			want: nil,
		},
		{
			name: "invalid authentication",
			args: args{
				cs: "TLS_ECDHE_RSA_WITHOUT_AES_256_GCM_SHA384",
			},
			want: nil,
		},
		{
			name: "invalid session cipher algorithm",
			args: args{
				cs: "TLS_ECDHE_RSA_WITH_AIS_256_GCM_SHA384",
			},
			want: nil,
		},
		{
			name: "invalid session cipher key length",
			args: args{
				cs: "TLS_ECDHE_RSA_WITH_AES_257_GCM_SHA384",
			},
			want: nil,
		},
		{
			name: "invalid session cipher mode",
			args: args{
				cs: "TLS_ECDHE_RSA_WITH_AES_256_FCM_SHA384",
			},
			want: nil,
		},
		{
			name: "invalid mac algorithm",
			args: args{
				cs: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHO384",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tlsCipherSuites(tt.args.cs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_tlsVersion(t *testing.T) {
	type args struct {
		version *string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{
			name: "1_3",
			args: args{
				version: util.Ref("1_3"),
			},
			want: 1.3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tlsVersion(tt.args.version); got != tt.want {
				t.Errorf("tlsVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_azureStorageDiscovery_discoverDiagnosticSettings(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		resourceURI string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ontology.ActivityLogging
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "No Diagnostic Setting available",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
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
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				resourceURI: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
			},
			want: &ontology.ActivityLogging{
				Enabled:           true,
				LoggingServiceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/insights-integration/providers/Microsoft.OperationalInsights/workspaces/workspace1"},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Init Diagnostic Settings Client
			_ = d.initDiagnosticsSettingsClient()

			got, raw, err := d.discoverDiagnosticSettings(tt.args.resourceURI)

			tt.wantErr(t, err)
			if tt.wantErr != nil {
				assert.NotNil(t, raw)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
