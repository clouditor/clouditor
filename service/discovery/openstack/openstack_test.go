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
	"fmt"
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
					WithTargetOfEvaluationID(testdata.MockTargetOfEvaluationID2),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				return assert.Equal(t, "OpenStack", got.Name())
			},
		},
		{
			name: "Happy path: with target of evaluation id",
			args: args{
				opts: []DiscoveryOption{
					WithAuthorizer(gophercloud.AuthOptions{
						IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
						Username:         testdata.MockOpenstackUsername,
						Password:         testdata.MockOpenstackPassword,
						TenantName:       testdata.MockOpenstackTenantName,
						AllowReauth:      true,
					}),
					WithTargetOfEvaluationID(testdata.MockTargetOfEvaluationID2),
				},
			},
			want: func(t *testing.T, got discovery.Discoverer) bool {
				assert.Equal(t, testdata.MockTargetOfEvaluationID2, got.TargetOfEvaluationID())
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
				assert.Equal(t, config.DefaultTargetOfEvaluationID, got.TargetOfEvaluationID())
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
			name: "cluster client error",
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
							if eo.Type == "container-infrastructure-management" {
								return "", errors.New("this is a test error")
							}
							return testhelper.Endpoint(), nil
						},
					},
				},
			},
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not create cluster client:")
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
				assert.True(t, got.AllowReauth)
				got.AllowReauth = false // We do not want to check this field in the following
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
						envVariableValue: testdata.MockOpenstackProjectID1,
					},
					{
						envVariableKey:   "OS_PROJECT_ID",
						envVariableValue: testdata.MockOpenstackProjectID1,
					},
				},
			},
			want: func(t *testing.T, got gophercloud.AuthOptions) bool {
				want := gophercloud.AuthOptions{
					IdentityEndpoint: testdata.MockOpenstackIdentityEndpoint,
					Username:         testdata.MockOpenstackUsername,
					Password:         testdata.MockOpenstackPassword,
					TenantID:         testdata.MockOpenstackProjectID1,
					AllowReauth:      true,
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

	type fields struct {
		ctID     string
		clients  clients
		authOpts *gophercloud.AuthOptions
		domain   *domain
		project  *project
		projects map[string]ontology.IsResource

		testhelper string
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error authorization",
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not authorize openstack:")
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
				project: &project{},
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				projects: map[string]ontology.IsResource{},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover servers:")
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
				project: &project{
					projectID:   testdata.MockOpenstackProjectID1,
					projectName: testdata.MockOpenstackProjectName1,
				},
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				projects: map[string]ontology.IsResource{},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover network interfaces:")
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
				project: &project{
					projectID:   testdata.MockOpenstackProjectID1,
					projectName: testdata.MockOpenstackProjectName1,
				},
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				projects: map[string]ontology.IsResource{},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover block storage:")
			},
		},
		{
			name: "error discover clusters",
			fields: fields{
				testhelper: "clusters",
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
				project: &project{
					projectID:   testdata.MockOpenstackProjectID1,
					projectName: testdata.MockOpenstackProjectName1,
				},
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				projects: map[string]ontology.IsResource{},
			},
			want: assert.Nil[[]ontology.IsResource],
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not discover clusters:")
			},
		},
		// {
		// name: "error discover projects",
		// Not possible to test. The method discoverServer() is called before discoverProjects() and adds the project ID which is the only possibility to get an error.
		// },
		{
			name: "error discover projects: but there is no error, as a resource is added based on other information discovered before.",
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
				project: &project{},
				domain: &domain{
					domainID: "test domain ID",
				},
				projects: map[string]ontology.IsResource{},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				_, ok := got[10].(*ontology.ResourceGroup)
				assert.True(t, ok)
				_, ok = got[11].(*ontology.ResourceGroup)
				assert.True(t, ok)
				_, ok = got[12].(*ontology.ResourceGroup)
				assert.True(t, ok)
				_, ok = got[13].(*ontology.ResourceGroup)
				return assert.True(t, ok)
			},
			wantErr: assert.NoError,
		},
		// {
		// 	name: "error discover domains",
		// We can’t test this, because the domain is still added when domainID is unset—and in that case, server discovery fails which is done before.
		// },
		{
			name: "Happy path: with one project in map that is nil",
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
				project: &project{
					projectID:   testdata.MockOpenstackProjectID1,
					projectName: testdata.MockOpenstackProjectName1,
				},
				domain: &domain{
					domainID: testdata.MockOpenstackDomainID1,
				},
				projects: map[string]ontology.IsResource{
					testdata.MockOpenstackProjectID1: nil,
				},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				return assert.Equal(t, 14, len(got))
			},
			wantErr: assert.NoError,
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
				project: &project{},
				domain: &domain{
					domainID: "test domain ID",
				},
				projects: map[string]ontology.IsResource{},
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				return assert.Equal(t, 14, len(got))
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testhelper.SetupHTTP()

			d := &openstackDiscovery{
				ctID:               tt.fields.ctID,
				clients:            tt.fields.clients,
				authOpts:           tt.fields.authOpts,
				domain:             tt.fields.domain,
				configuredProject:  tt.fields.project,
				discoveredProjects: tt.fields.projects,
			}

			switch tt.fields.testhelper {
			case "all":
				fmt.Println("Setting up handlers for all resources")
				const ConsoleOutputBody = `{
					"output": "output test"
				}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
				openstacktest.HandleNetworkListSuccessfully(t)
				openstacktest.MockStorageListResponse(t)
				openstacktest.HandleListClusterSuccessfully(t)
			case "domain":
				fmt.Println("Setting up handlers to get an error for domain resources")
				const ConsoleOutputBody = `{
					"output": "output test"
				}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
				openstacktest.HandleNetworkListSuccessfully(t)
				openstacktest.MockStorageListResponse(t)
				openstacktest.HandleListClusterSuccessfully(t)
			case "project":
				fmt.Println("Setting up handlers to get an error for project resources")
				const ConsoleOutputBody = `{
					"output": "output test"
				}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
				openstacktest.HandleNetworkListSuccessfully(t)
				openstacktest.MockStorageListResponse(t)
				openstacktest.HandleListClusterSuccessfully(t)
			case "clusters":
				fmt.Println("Setting up handlers to get an error for storage resources")
				const ConsoleOutputBody = `{
					"output": "output test"
				}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
				openstacktest.HandleNetworkListSuccessfully(t)
				openstacktest.MockStorageListResponse(t)
			case "storage":
				fmt.Println("Setting up handlers to get an error for storage resources")
				const ConsoleOutputBody = `{
					"output": "output test"
				}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
				openstacktest.HandleNetworkListSuccessfully(t)
			case "network":
				fmt.Println("Setting up handlers to get an error for network resources")
				const ConsoleOutputBody = `{
						"output": "output test"
					}`

				openstacktest.HandleServerListSuccessfully(t)
				openstacktest.HandleShowConsoleOutputSuccessfully(t, ConsoleOutputBody)
				openstacktest.HandleInterfaceListSuccessfully(t)
			}

			gotList, err := d.List()

			tt.want(t, gotList)
			tt.wantErr(t, err)
			testhelper.TeardownHTTP()
		})
	}
}
