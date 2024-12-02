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

package openstack

import (
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"github.com/gophercloud/gophercloud/v2"
)

func TestNewOpenstackDiscovery(t *testing.T) {
	type args struct {
		opts []DiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want assert.Want[discovery.Discoverer]
	}{
		{
			name: "error: oauthOpts not set",
			args: args{},
			want: assert.Nil[discovery.Discoverer],
		},
		{
			name: "Happy path: Name",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(gophercloud.AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
					WithCertificationTargetID(testdata.MockCertificationTargetID2),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				return assert.Equal(t, "OpenStack", got.Name())
			},
		},
		{
			name: "Happy path: with certification target id",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(gophercloud.AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
					WithCertificationTargetID(testdata.MockCertificationTargetID2),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				assert.Equal(t, testdata.MockCertificationTargetID2, got.CertificationTargetID())
				return assert.NotNil(t, got)
			},
		},
		{
			name: "Happy path: with authorizer",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(gophercloud.AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				assert.Equal(t, config.DefaultCertificationTargetID, got.CertificationTargetID())
				return assert.NotNil(t, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOpenstackDiscovery(tt.args.opts...)
			tt.want(t, got)
		})
	}
}

// func Test_openstackDiscovery_authorize(t *testing.T) {
// 	testhelper.SetupHTTP()
// 	defer testhelper.TeardownHTTP()

// 	type fields struct {
// 		csID     string
// 		provider *gophercloud.ProviderClient
// 		clients  clients
// 		authOpts *gophercloud.AuthOptions
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		wantErr assert.ErrorAssertionFunc
// 	}{
// 		{
// 			name: "Happy path: already authorized",
// 			fields: fields{
// 				provider: &gophercloud.ProviderClient{},
// 			},
// 			wantErr: assert.NoError,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			d := &openstackDiscovery{
// 				csID:     tt.fields.csID,
// 				clients:  tt.fields.clients,
// 				authOpts: tt.fields.authOpts,
// 			}
// 			tt.wantErr(t, d.authorize())
// 		})
// 	}
// }

func TestNewAuthorizer(t *testing.T) {
	type envVariables struct {
		envVariableKey   string
		envVariableValue string
	}
	type fields struct {
		envVariables []envVariables
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[gophercloud.AuthOptions]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				envVariables: []envVariables{
					{
						envVariableKey:   "OS_AUTH_URL",
						envVariableValue: testdata.MockOpenstackIdentityEndpoint,
					},
					{
						envVariableKey:   "OS_USERNAME",
						envVariableValue: testdata.MockOpenstackUsername,
					},
					{
						envVariableKey:   "OS_PASSWORD",
						envVariableValue: testdata.MockOpenstackPassword,
					},
					{
						envVariableKey:   "OS_TENANT_ID",
						envVariableValue: testdata.MockProjectID1,
					},
					{
						envVariableKey:   "OS_PROJECT_ID",
						envVariableValue: testdata.MockProjectID1,
					},
					{
						envVariableKey:   "OS_SYSTEM_SCOPE",
						envVariableValue: "true",
					},
				},
			},
			want: func(t *testing.T, got gophercloud.AuthOptions) bool {
				want := gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantID:         testdata.MockProjectID1,
					AllowReauth:      true,
				}
				return assert.Equal(t, want, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env variables
			for _, env := range tt.fields.envVariables {
				if env.envVariableKey != "" {
					t.Setenv(env.envVariableKey, env.envVariableValue)
				}
			}
			got, err := NewAuthorizer()

			tt.want(t, got)
			tt.wantErr(t, err)
		})
	}
}
