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
	"errors"
	"testing"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/testutil/servicetest/discoverytest/openstacktest"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
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

func Test_openstackDiscovery_authorize(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "error authentication",
			fields: fields{},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error while authenticating:")
			},
		},
		{
			name: "compute client error",
			fields: fields{
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return "", errors.New("this is a test error")
						},
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not create compute client:")
			},
		},
		{
			name: "network client error",
			fields: fields{
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							if eo.Type == "network" {
								return "", errors.New("this is a test error")
							}
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not create network client:")
			},
		},
		{
			name: "storage client error",
			fields: fields{
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							if eo.Type == "block-storage" {
								return "", errors.New("this is a test error")
							}
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not create block storage client:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				ctID:     tt.fields.ctID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}

			err := d.authorize()

			tt.wantErr(t, err)
		})
	}
}

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
			name: "error: missing OS_AUTH_URL",
			fields: fields{
				envVariables: []envVariables{},
			},
			want: func(t *testing.T, got gophercloud.AuthOptions) bool {
				return assert.Empty(t, got)
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "Missing environment variable [OS_AUTH_URL]")
			},
		},
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
				},
			},
			want: func(t *testing.T, got gophercloud.AuthOptions) bool {
				want := gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantID:         testdata.MockProjectID1,
				}
				return assert.Equal(t, want, got)
			},
			wantErr: assert.NoError,
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

func Test_openstackDiscovery_List(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	type fields struct {
		ctID       string
		clients    clients
		authOpts   *gophercloud.AuthOptions
		testhelper string
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error discover domains",
			fields: fields{
				testhelper: "",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover domains:")
			},
		},
		{
			name: "error discover projects",
			fields: fields{
				testhelper: "project",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover projects:")
			},
		},
		{
			name: "error discover network interfaces",
			fields: fields{
				testhelper: "network",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover network interfaces:")
			},
		},
		{
			name: "error discover server",
			fields: fields{
				testhelper: "server",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover servers:")
			},
		},
		{
			name: "error discover block storage",
			fields: fields{
				testhelper: "storage",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover block storage:")
			},
		},
		{
			name: "Happy path",
			fields: fields{
				testhelper: "all",
				authOpts: &gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantName:       testdata.MockOpenstackTenantName,
				},
				clients: clients{
					provider: &gophercloud.ProviderClient{
						TokenID: client.TokenID,
						EndpointLocator: func(eo gophercloud.EndpointOpts) (string, error) {
							return testhelper.Endpoint(), nil
						},
					},
					identityClient: client.ServiceClient(),
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				return assert.Equal(t, 11, len(got))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &openstackDiscovery{
				ctID:     tt.fields.ctID,
				clients:  tt.fields.clients,
				authOpts: tt.fields.authOpts,
			}

			// The ordering is important, otherwise there are errors if a handleX function is called twice.
			switch tt.fields.testhelper {
			case "all":
				openstacktest.MockStorageListResponse(t)
			case "project":
				openstacktest.HandleListDomainsSuccessfully(t)
			case "network":
				openstacktest.HandleListProjectsSuccessfully(t)
			case "server":
				openstacktest.HandleNetworkListSuccessfully(t)
			case "storage":
				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleInterfaceListSuccessfully(t)
			}

			gotList, err := d.List()

			tt.want(t, gotList)
			tt.wantErr(t, err)
		})
	}
}
