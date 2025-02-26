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
	"clouditor.io/clouditor/v2/internal/testutil/assert"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_azureComputeDiscovery_discoverFunctionsWebApps(t *testing.T) {
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
			name: "Happy path",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []ontology.IsResource{
				&ontology.Function{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/function1",
					Name:         "function1",
					CreationTime: nil,
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					GeoLocation: &ontology.GeoLocation{
						Region: "West Europe",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1\",\"kind\":\"functionapp,linux\",\"location\":\"West Europe\",\"name\":\"function1\",\"properties\":{\"publicNetworkAccess\":\"Enabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{\"linuxFxVersion\":\"PYTHON|3.8\"}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"function1\",\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.1\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
					NetworkInterfaceIds: []string{},
					ResourceLogging: &ontology.ResourceLogging{
						Enabled: false,
					},
					/*HttpEndpoint: &ontology.HttpEndpoint{
						TransportEncryption: &ontology.TransportEncryption{
							Enabled:    true,
							Enforced:   false,
							TlsVersion: constants.TLS1_1,
							Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},*/
					InternetAccessibleEndpoint: true,
					Redundancies:               []*ontology.Redundancy{},
					RuntimeVersion:             "3.8",
					RuntimeLanguage:            "PYTHON",
				},
				&ontology.Function{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/function2",
					Name:         "function2",
					CreationTime: nil,
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					GeoLocation: &ontology.GeoLocation{
						Region: "West Europe",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2\",\"kind\":\"functionapp\",\"location\":\"West Europe\",\"name\":\"function2\",\"properties\":{\"publicNetworkAccess\":\"Disabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"function2\",\"properties\":{\"javaVersion\":\"1.8\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
					NetworkInterfaceIds: []string{},
					ResourceLogging: &ontology.ResourceLogging{
						Enabled: false,
					},
					/*HttpEndpoint: &ontology.HttpEndpoint{
						TransportEncryption: &ontology.TransportEncryption{
							Enabled:    false,
							Enforced:   false,
							TlsVersion: "",
							Algorithm:  "",
						},
					},*/
					InternetAccessibleEndpoint: false,
					Redundancies:               []*ontology.Redundancy{},
					RuntimeVersion:             "1.8",
					RuntimeLanguage:            "Java",
				},
				&ontology.Function{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/webapp1",
					Name:         "WebApp1",
					CreationTime: nil,
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					GeoLocation: &ontology.GeoLocation{
						Region: "West Europe",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp1\",\"properties\":{\"httpsOnly\":true,\"publicNetworkAccess\":\"Enabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"WebApp1\",\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.1\"},\"type\":\"Microsoft.Web/sites/config\"}]}",
					NetworkInterfaceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/virtualnetworks/vnet1/subnets/subnet1"},
					ResourceLogging: &ontology.ResourceLogging{
						Enabled: true,
					},
					/*HttpEndpoint: &ontology.HttpEndpoint{
						TransportEncryption: &ontology.TransportEncryption{
							Enabled:    true,
							Enforced:   true,
							TlsVersion: constants.TLS1_1,
							Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},*/
					InternetAccessibleEndpoint: true,
					Redundancies:               []*ontology.Redundancy{},
				},
				&ontology.Function{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/webapp2",
					Name:         "WebApp2",
					CreationTime: nil,
					Labels: map[string]string{
						"testKey1": "testTag1",
						"testKey2": "testTag2",
					},
					GeoLocation: &ontology.GeoLocation{
						Region: "West Europe",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2\",\"kind\":\"app,linux\",\"location\":\"West Europe\",\"name\":\"WebApp2\",\"properties\":{\"httpsOnly\":false,\"publicNetworkAccess\":\"Disabled\",\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"name\":\"WebApp2\",\"properties\":{},\"type\":\"Microsoft.Web/sites/config\"}]}",
					NetworkInterfaceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/virtualnetworks/vnet1/subnets/subnet2"},
					ResourceLogging: &ontology.ResourceLogging{
						Enabled: false,
					},
					/*HttpEndpoint: &ontology.HttpEndpoint{
						TransportEncryption: &ontology.TransportEncryption{
							Enabled:    false,
							Enforced:   false,
							TlsVersion: "",
							Algorithm:  "",
						},
					},*/
					InternetAccessibleEndpoint: false,
					Redundancies:               []*ontology.Redundancy{},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverFunctionsWebApps()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleFunction(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
		clientWebApps  bool
	}
	type args struct {
		function *armappservice.Site
		config   armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ontology.IsResource
	}{
		{
			name: "Empty input",
			args: args{
				function: nil,
			},
			want: nil,
		},
		{
			name: "Happy path: Linux function",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				function: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1"),
					Name:     util.Ref("function1"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("functionapp,linux"),
					Properties: &armappservice.SiteProperties{
						SiteConfig: &armappservice.SiteConfig{
							LinuxFxVersion: util.Ref("PYTHON|3.8"),
						},
						HTTPSOnly:     util.Ref(true),
						ResourceGroup: util.Ref("res1"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &ontology.Function{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/function1",
				Name:         "function1",
				CreationTime: nil,
				Labels: map[string]string{
					"testKey1": "testTag1",
					"testKey2": "testTag2",
				},
				GeoLocation: &ontology.GeoLocation{
					Region: "West Europe",
				},
				ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function1\",\"kind\":\"functionapp,linux\",\"location\":\"West Europe\",\"name\":\"function1\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{\"linuxFxVersion\":\"PYTHON|3.8\"}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
				NetworkInterfaceIds: []string{},
				ResourceLogging: &ontology.ResourceLogging{
					Enabled: false,
				},
				InternetAccessibleEndpoint: false,
				Redundancies:               []*ontology.Redundancy{},
				RuntimeVersion:             "3.8",
				RuntimeLanguage:            "PYTHON",
				/*HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},
				PublicAccess:    false,*/
			},
		},
		{
			name: "Happy path: Windows function",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				function: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2"),
					Name:     util.Ref("function2"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("functionapp"),
					Properties: &armappservice.SiteProperties{
						SiteConfig:    &armappservice.SiteConfig{},
						ResourceGroup: util.Ref("res1"),
						HTTPSOnly:     util.Ref(true),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							JavaVersion:       util.Ref("1.8"),
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &ontology.Function{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/function2",
				Name:         "function2",
				CreationTime: nil,
				Labels: map[string]string{
					"testKey1": "testTag1",
					"testKey2": "testTag2",
				},
				GeoLocation: &ontology.GeoLocation{
					Region: "West Europe",
				},
				ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/function2\",\"kind\":\"functionapp\",\"location\":\"West Europe\",\"name\":\"function2\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{}},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"javaVersion\":\"1.8\",\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
				NetworkInterfaceIds: []string{},
				ResourceLogging: &ontology.ResourceLogging{
					Enabled: false,
				},
				InternetAccessibleEndpoint: false,
				Redundancies:               []*ontology.Redundancy{},
				RuntimeVersion:             "1.8",
				RuntimeLanguage:            "Java",
				/*HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enforced:   true,
						Enabled:    true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},
				*/
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Set clients if needed
			if tt.fields.clientWebApps {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			assert.Equal(t, tt.want, d.handleFunction(tt.args.function, tt.args.config))
		})
	}
}

func Test_azureComputeDiscovery_discoverVirtualMachines(t *testing.T) {
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
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: []ontology.IsResource{
				&ontology.VirtualMachine{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/virtualmachines/vm1",
					Name:         "vm1",
					CreationTime: timestamppb.New(creationTime),
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"location\":\"eastus\",\"name\":\"vm1\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true,\"storageUri\":\"https://logstoragevm1.blob.core.windows.net/\"}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"123\"},{\"id\":\"234\"}]},\"osProfile\":{\"linuxConfiguration\":{\"patchSettings\":{\"patchMode\":\"AutomaticByPlatform\"}}},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_1\"}},{\"managedDisk\":{\"id\":\"data_disk_2\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"},\"resources\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1/extensions/MicrosoftMonitoringAgent\"}]}]}",
					NetworkInterfaceIds: []string{"123", "234"},
					BlockStorageIds:     []string{"os_test_disk", "data_disk_1", "data_disk_2"},
					BootLogging: &ontology.BootLogging{
						Enabled: true,
						//LoggingService: []ontology.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
						LoggingServiceIds: []string{},
						RetentionPeriod:   durationpb.New(0),
					},
					OsLogging: &ontology.OSLogging{
						Enabled:           true,
						RetentionPeriod:   durationpb.New(0),
						LoggingServiceIds: []string{},
					},
					AutomaticUpdates: &ontology.AutomaticUpdates{
						Enabled:  true,
						Interval: durationpb.New(Duration30Days),
					},
					MalwareProtection: &ontology.MalwareProtection{},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled:           true,
						RetentionPeriod:   durationpb.New(RetentionPeriod90Days),
						LoggingServiceIds: []string{},
					},
				},
				&ontology.VirtualMachine{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/virtualmachines/vm2",
					Name:         "vm2",
					CreationTime: nil,
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2\",\"location\":\"eastus\",\"name\":\"vm2\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"987\"},{\"id\":\"654\"}]},\"osProfile\":{\"windowsConfiguration\":{\"enableAutomaticUpdates\":true,\"patchSettings\":{\"patchMode\":\"AutomaticByOS\"}}},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_2\"}},{\"managedDisk\":{\"id\":\"data_disk_3\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}}},\"resources\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm2/extensions/OmsAgentForLinux\"}]}]}",
					NetworkInterfaceIds: []string{"987", "654"},
					BlockStorageIds:     []string{"os_test_disk", "data_disk_2", "data_disk_3"},
					BootLogging: &ontology.BootLogging{
						Enabled:           true,
						LoggingServiceIds: []string{},
						RetentionPeriod:   durationpb.New(0),
					},
					OsLogging: &ontology.OSLogging{
						Enabled:           true,
						LoggingServiceIds: []string{},
						RetentionPeriod:   durationpb.New(0),
					},
					AutomaticUpdates: &ontology.AutomaticUpdates{
						Enabled:  true,
						Interval: durationpb.New(Duration30Days),
					},
					MalwareProtection: &ontology.MalwareProtection{},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled:           true,
						RetentionPeriod:   durationpb.New(RetentionPeriod90Days),
						LoggingServiceIds: []string{},
					},
				},
				&ontology.VirtualMachine{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/virtualmachines/vm3",
					Name:         "vm3",
					CreationTime: nil,
					Labels:       map[string]string{},
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:                 "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm3\",\"location\":\"eastus\",\"name\":\"vm3\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{}}}}]}",
					NetworkInterfaceIds: []string{},
					BlockStorageIds:     []string{},
					BootLogging: &ontology.BootLogging{
						Enabled:           false,
						LoggingServiceIds: []string{},
						RetentionPeriod:   durationpb.New(0),
					},
					OsLogging: &ontology.OSLogging{
						Enabled:           false,
						LoggingServiceIds: []string{},
						RetentionPeriod:   durationpb.New(0),
					},
					AutomaticUpdates: &ontology.AutomaticUpdates{
						Enabled:  false,
						Interval: nil,
					},
					MalwareProtection: &ontology.MalwareProtection{},
					ActivityLogging: &ontology.ActivityLogging{
						Enabled:           true,
						RetentionPeriod:   durationpb.New(RetentionPeriod90Days),
						LoggingServiceIds: []string{},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.discoverVirtualMachines()
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleVirtualMachines(t *testing.T) {
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ontology.IsResource
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Virtual Machine is empty",
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrEmptyVirtualMachine)
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender(), WithDefenderProperties(map[string]*defenderProperties{
					DefenderVirtualMachineType: {
						monitoringLogDataEnabled: true,
						securityAlertsEnabled:    true,
					},
				})),
			},
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"),
					Name:     util.Ref("vm1"),
					Location: util.Ref("eastus"),
					Properties: &armcompute.VirtualMachineProperties{
						TimeCreated: util.Ref(creationTime),
						NetworkProfile: &armcompute.NetworkProfile{
							NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
								{
									ID: util.Ref("123"),
								},
								{
									ID: util.Ref("234"),
								},
							},
						},
						StorageProfile: &armcompute.StorageProfile{
							OSDisk: &armcompute.OSDisk{
								ManagedDisk: &armcompute.ManagedDiskParameters{
									ID: util.Ref("os_test_disk"),
								},
							},
							DataDisks: []*armcompute.DataDisk{
								{
									ManagedDisk: &armcompute.ManagedDiskParameters{
										ID: util.Ref("data_disk_1"),
									},
								},
								{
									ManagedDisk: &armcompute.ManagedDiskParameters{
										ID: util.Ref("data_disk_2"),
									},
								},
							},
						},
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    util.Ref(true),
								StorageURI: util.Ref("https://logstoragevm1.blob.core.windows.net/"),
							},
						},
					},
				},
			},
			want: &ontology.VirtualMachine{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/virtualmachines/vm1",
				Name:         "vm1",
				CreationTime: timestamppb.New(creationTime),
				Labels:       map[string]string{},
				GeoLocation: &ontology.GeoLocation{
					Region: "eastus",
				},
				ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:                 "{\"*armcompute.VirtualMachine\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"location\":\"eastus\",\"name\":\"vm1\",\"properties\":{\"diagnosticsProfile\":{\"bootDiagnostics\":{\"enabled\":true,\"storageUri\":\"https://logstoragevm1.blob.core.windows.net/\"}},\"networkProfile\":{\"networkInterfaces\":[{\"id\":\"123\"},{\"id\":\"234\"}]},\"storageProfile\":{\"dataDisks\":[{\"managedDisk\":{\"id\":\"data_disk_1\"}},{\"managedDisk\":{\"id\":\"data_disk_2\"}}],\"osDisk\":{\"managedDisk\":{\"id\":\"os_test_disk\"}}},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"}}]}",
				NetworkInterfaceIds: []string{"123", "234"},
				BlockStorageIds:     []string{"os_test_disk", "data_disk_1", "data_disk_2"},
				BootLogging: &ontology.BootLogging{
					Enabled: true,
					//LoggingService: []ontology.ResourceID{"https://logstoragevm1.blob.core.windows.net/"},
					LoggingServiceIds:        []string{},
					RetentionPeriod:          durationpb.New(0),
					MonitoringLogDataEnabled: true,
					SecurityAlertsEnabled:    true,
				},
				OsLogging: &ontology.OSLogging{
					Enabled:                  false,
					LoggingServiceIds:        []string{},
					RetentionPeriod:          durationpb.New(0),
					MonitoringLogDataEnabled: true,
					SecurityAlertsEnabled:    true,
				},
				AutomaticUpdates: &ontology.AutomaticUpdates{
					Enabled:  false,
					Interval: nil,
				},
				MalwareProtection: &ontology.MalwareProtection{},
				ActivityLogging: &ontology.ActivityLogging{
					Enabled:           true,
					RetentionPeriod:   durationpb.New(RetentionPeriod90Days),
					LoggingServiceIds: []string{},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			got, err := d.handleVirtualMachines(tt.args.vm)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isBootDiagnosticEnabled(t *testing.T) {
	ID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"
	name := "vm1"
	enabledTrue := true

	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty input",
			args: args{
				vm: nil,
			},
			want: false,
		},
		{
			name: "Empty properties value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:         &ID,
					Name:       &name,
					Properties: nil,
				},
			},
			want: false,
		},
		{
			name: "Empty DiagnosticsProfile value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: nil,
					},
				},
			},
			want: false,
		},
		{
			name: "Empty BootDiagnostics value",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: nil,
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Correct input",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled: &enabledTrue,
							},
						},
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isBootDiagnosticEnabled(tt.args.vm))
		})
	}
}

func Test_bootLogOutput(t *testing.T) {
	ID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1"
	name := "vm1"
	enabledTrue := true
	enabledFalse := false
	storageUri := "https://logstoragevm1.blob.core.windows.net/"
	emptyStorageUri := ""

	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty input",
			args: args{
				vm: nil,
			},
			want: "",
		},
		{
			name: "StorageURI is nil",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledFalse,
								StorageURI: nil,
							},
						},
					},
				},
			},
			want: "",
		},
		{
			name: "BootDiagnostics disabled",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledFalse,
								StorageURI: &emptyStorageUri,
							},
						},
					},
				},
			},
			want: "",
		},
		{
			name: "BootDiagnostics enabled",
			args: args{
				vm: &armcompute.VirtualMachine{
					ID:   &ID,
					Name: &name,
					Properties: &armcompute.VirtualMachineProperties{
						DiagnosticsProfile: &armcompute.DiagnosticsProfile{
							BootDiagnostics: &armcompute.BootDiagnostics{
								Enabled:    &enabledTrue,
								StorageURI: &storageUri,
							},
						},
					},
				},
			},
			//want: storageUri,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, bootLogOutput(tt.args.vm))
		})
	}
}

func Test_azureComputeDiscovery_discoverBlockStorage(t *testing.T) {
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
			want: []ontology.IsResource{
				&ontology.BlockStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/disks/disk1",
					Name:         "disk1",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1\",\"location\":\"eastus\",\"managedBy\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"name\":\"disk1\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"\",\"type\":\"EncryptionAtRestWithPlatformKey\"},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[null]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
						},
					},
				},
				&ontology.BlockStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/disks/disk2",
					Name:         "disk2",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk2\",\"location\":\"eastus\",\"managedBy\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"name\":\"disk2\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1\",\"type\":\"EncryptionAtRestWithCustomerKey\"},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault1\",\"location\":\"germanywestcentral\",\"name\":\"encryptionkeyvault1\",\"properties\":{\"activeKey\":{\"keyUrl\":\"https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382\",\"sourceVault\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault1\"}}},\"type\":\"Microsoft.Compute/diskEncryptionSets\"}]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
							CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
								Algorithm: "",
								Enabled:   true,
								KeyUrl:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
							},
						},
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
						},
					},
				},
				&ontology.BlockStorage{
					Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2/providers/microsoft.compute/disks/disk3",
					Name:         "disk3",
					CreationTime: timestamppb.New(creationTime),
					GeoLocation: &ontology.GeoLocation{
						Region: "eastus",
					},
					Labels:   map[string]string{},
					ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
					Raw:      "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res2/providers/Microsoft.Compute/disks/disk3\",\"location\":\"eastus\",\"managedBy\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1\",\"name\":\"disk3\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"\",\"type\":\"EncryptionAtRestWithPlatformKey\"},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"},\"type\":\"Microsoft.Compute/disks\"}],\"*armcompute.DiskEncryptionSet\":[null]}",
					AtRestEncryption: &ontology.AtRestEncryption{
						Type: &ontology.AtRestEncryption_ManagedKeyEncryption{
							ManagedKeyEncryption: &ontology.ManagedKeyEncryption{
								Algorithm: "AES256",
								Enabled:   true,
							},
						},
					},
					Backups: []*ontology.Backup{
						{
							Enabled:         false,
							RetentionPeriod: nil,
							Interval:        nil,
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

			got, err := d.discoverBlockStorages()
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_handleBlockStorage(t *testing.T) {
	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		disk *armcompute.Disk
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ontology.BlockStorage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			args: args{
				disk: nil,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "disk is nil")
			},
		},
		{
			name: "Empty diskID",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID: &diskID,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type: &encType,
						},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get block storage properties for the atRestEncryption:")
			},
		},
		{
			name: "Empty encryptionType",
			args: args{
				disk: &armcompute.Disk{
					ID: &diskID,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type: nil,
						},
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting atRestEncryption properties of blockStorage")
			},
		},
		{
			name: "No error",
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type:                &encType,
							DiskEncryptionSetID: &encSetID,
						},
						TimeCreated: &creationTime,
					},
					ManagedBy: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2/providers/microsoft.compute/disks/disk3"),
				},
			},
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want: &ontology.BlockStorage{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.compute/disks/disk1",
				Name:         "disk1",
				CreationTime: timestamppb.New(creationTime),
				GeoLocation: &ontology.GeoLocation{
					Region: "eastus",
				},
				Labels:   map[string]string{},
				ParentId: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2"),
				Raw:      "{\"*armcompute.Disk\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1\",\"location\":\"eastus\",\"managedBy\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res2/providers/microsoft.compute/disks/disk3\",\"name\":\"disk1\",\"properties\":{\"encryption\":{\"diskEncryptionSetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1\",\"type\":\"EncryptionAtRestWithCustomerKey\"},\"timeCreated\":\"2017-05-24T13:28:53.004540398Z\"}}],\"*armcompute.DiskEncryptionSet\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryption-keyvault1\",\"location\":\"germanywestcentral\",\"name\":\"encryptionkeyvault1\",\"properties\":{\"activeKey\":{\"keyUrl\":\"https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382\",\"sourceVault\":{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.KeyVault/vaults/keyvault1\"}}},\"type\":\"Microsoft.Compute/diskEncryptionSets\"}]}",
				AtRestEncryption: &ontology.AtRestEncryption{
					Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
						CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
							Algorithm: "",
							Enabled:   true,
							KeyUrl:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
						},
					},
				},
				Backups: []*ontology.Backup{
					{
						Enabled:         false,
						RetentionPeriod: nil,
						Interval:        nil,
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, err := d.handleBlockStorage(tt.args.disk)
			if !tt.wantErr(t, err, fmt.Sprintf("handleBlockStorage(%v)", tt.args.disk)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_blockStorageAtRestEncryption(t *testing.T) {
	encType := armcompute.EncryptionTypeEncryptionAtRestWithCustomerKey
	diskID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/disks/disk1"
	diskName := "disk1"
	diskRegion := "eastus"
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	creationTime := time.Date(2017, 05, 24, 13, 28, 53, 4540398, time.UTC)

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		disk *armcompute.Disk
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ontology.AtRestEncryption
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty disk",
			args: args{},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "disk is empty")
			},
		},
		{
			name: "Error getting atRestEncryptionProperties",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption:  &armcompute.Encryption{},
						TimeCreated: &creationTime,
					},
				},
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "error getting atRestEncryption properties of blockStorage")
			},
		},
		{
			name: "No error",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				disk: &armcompute.Disk{
					ID:       &diskID,
					Name:     &diskName,
					Location: &diskRegion,
					Properties: &armcompute.DiskProperties{
						Encryption: &armcompute.Encryption{
							Type:                &encType,
							DiskEncryptionSetID: &encSetID,
						},
						TimeCreated: &creationTime,
					},
				},
			},
			want: &ontology.AtRestEncryption{
				Type: &ontology.AtRestEncryption_CustomerKeyEncryption{
					CustomerKeyEncryption: &ontology.CustomerKeyEncryption{
						Algorithm: "",
						Enabled:   true,
						KeyUrl:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, _, err := d.blockStorageAtRestEncryption(tt.args.disk)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_azureComputeDiscovery_keyURL(t *testing.T) {
	encSetID := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1"
	encSetID2 := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault2"

	type fields struct {
		azureDiscovery *azureDiscovery
	}
	type args struct {
		diskEncryptionSetID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Empty input",
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrMissingDiskEncryptionSetID)
			},
		},
		{
			name:   "Error get disc encryption set",
			fields: fields{&azureDiscovery{}},
			args: args{
				diskEncryptionSetID: encSetID,
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get key vault:")
			},
		},
		{
			name: "Empty keyURL",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				diskEncryptionSetID: encSetID2,
			},
			want: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "could not get keyURL")
			},
		},
		{
			name: "No error",
			args: args{
				diskEncryptionSetID: encSetID,
			},
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			want:    "https://keyvault1.vault.azure.net/keys/customer-key/6273gdb374jz789hjm17819283748382",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			d := tt.fields.azureDiscovery

			got, _, err := d.keyURL(tt.args.diskEncryptionSetID)
			if !tt.wantErr(t, err, fmt.Sprintf("keyURL(%v)", tt.args.diskEncryptionSetID)) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_diskEncryptionSetName(t *testing.T) {
	type args struct {
		diskEncryptionSetID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Correct ID",
			args: args{
				diskEncryptionSetID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/diskEncryptionSets/encryptionkeyvault1",
			},
			want: "encryptionkeyvault1",
		},
		{
			name: "Empty ID",
			args: args{
				diskEncryptionSetID: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, diskEncryptionSetName(tt.args.diskEncryptionSetID))
		})
	}
}

func Test_runtimeInfo(t *testing.T) {
	type args struct {
		runtime string
	}
	tests := []struct {
		name                string
		args                args
		wantRuntimeLanguage string
		wantRuntimeVersion  string
	}{
		{
			name: "Empty input",
			args: args{
				runtime: "",
			},
			wantRuntimeLanguage: "",
			wantRuntimeVersion:  "",
		},
		{
			name: "Wrong input",
			args: args{
				runtime: "TEST",
			},
			wantRuntimeLanguage: "",
			wantRuntimeVersion:  "",
		},
		{
			name: "Happy path",
			args: args{
				runtime: "PYTHON|3.8",
			},
			wantRuntimeLanguage: "PYTHON",
			wantRuntimeVersion:  "3.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRuntimeLanguage, gotRuntimeVersion := runtimeInfo(tt.args.runtime)
			if gotRuntimeLanguage != tt.wantRuntimeLanguage {
				t.Errorf("runtimeInfo() gotRuntimeLanguage = %v, want %v", gotRuntimeLanguage, tt.wantRuntimeLanguage)
			}
			if gotRuntimeVersion != tt.wantRuntimeVersion {
				t.Errorf("runtimeInfo() gotRuntimeVersion = %v, want %v", gotRuntimeVersion, tt.wantRuntimeVersion)
			}
		})
	}
}

func Test_automaticUpdates(t *testing.T) {
	type args struct {
		vm *armcompute.VirtualMachine
	}
	tests := []struct {
		name                 string
		args                 args
		wantAutomaticUpdates *ontology.AutomaticUpdates
	}{
		{
			name:                 "Empty input",
			args:                 args{},
			wantAutomaticUpdates: &ontology.AutomaticUpdates{},
		},
		{
			name: "No automatic update for the given VM",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							LinuxConfiguration: &armcompute.LinuxConfiguration{
								PatchSettings: &armcompute.LinuxPatchSettings{},
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &ontology.AutomaticUpdates{},
		},
		{
			name: "Happy path: Linux",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							LinuxConfiguration: &armcompute.LinuxConfiguration{
								PatchSettings: &armcompute.LinuxPatchSettings{
									PatchMode: util.Ref(armcompute.LinuxVMGuestPatchModeAutomaticByPlatform),
								},
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &ontology.AutomaticUpdates{
				Enabled:  true,
				Interval: durationpb.New(Duration30Days),
			},
		},
		{
			name: "Happy path: Windows",
			args: args{
				vm: &armcompute.VirtualMachine{
					Properties: &armcompute.VirtualMachineProperties{
						OSProfile: &armcompute.OSProfile{
							WindowsConfiguration: &armcompute.WindowsConfiguration{
								PatchSettings: &armcompute.PatchSettings{
									PatchMode: util.Ref(armcompute.WindowsVMGuestPatchModeAutomaticByPlatform),
								},
								EnableAutomaticUpdates: util.Ref(true),
							},
						},
					},
				},
			},
			wantAutomaticUpdates: &ontology.AutomaticUpdates{
				Enabled:  true,
				Interval: durationpb.New(Duration30Days),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAutomaticUpdates := automaticUpdates(tt.args.vm)
			assert.Equal(t, tt.wantAutomaticUpdates, gotAutomaticUpdates)
		})
	}
}

func Test_azureComputeDiscovery_handleWebApp(t *testing.T) {
	type fields struct {
		azureDiscovery *azureDiscovery
		clientWebApps  bool
	}
	type args struct {
		webApp *armappservice.Site
		config armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ontology.IsResource
	}{
		{
			name: "Empty input",
			args: args{
				webApp: nil,
			},
			want: nil,
		},
		{
			name: "Happy path: WebApp Windows",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				webApp: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name:     util.Ref("WebApp1"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						HTTPSOnly:     util.Ref(true),
						ResourceGroup: util.Ref("res1"),
						SiteConfig: &armappservice.SiteConfig{
							MinTLSVersion: util.Ref(armappservice.SupportedTLSVersionsOne2),
						},
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			want: &ontology.Function{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/webapp1",
				Name:         "WebApp1",
				CreationTime: nil,
				Labels: map[string]string{
					"testKey1": "testTag1",
					"testKey2": "testTag2",
				},
				GeoLocation: &ontology.GeoLocation{
					Region: "West Europe",
				},
				ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp1\",\"properties\":{\"httpsOnly\":true,\"resourceGroup\":\"res1\",\"siteConfig\":{\"minTlsVersion\":\"1.2\"},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{\"minTlsCipherSuite\":\"TLS_AES_128_GCM_SHA256\",\"minTlsVersion\":\"1.2\"}}]}",
				NetworkInterfaceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/virtualnetworks/vnet1/subnets/subnet1"},
				ResourceLogging: &ontology.ResourceLogging{
					Enabled: true,
				},
				/*HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enabled:    true,
						Enforced:   true,
						TlsVersion: constants.TLS1_2,
						Algorithm:  string(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
					},
				},*/
				InternetAccessibleEndpoint: false,
				Redundancies:               []*ontology.Redundancy{},
			},
		},
		{
			name: "Happy path: WebApp Linux",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApps:  true,
			},
			args: args{
				webApp: &armappservice.Site{
					ID:       util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2"),
					Name:     util.Ref("WebApp2"),
					Location: util.Ref("West Europe"),
					Tags: map[string]*string{
						"testKey1": util.Ref("testTag1"),
						"testKey2": util.Ref("testTag2"),
					},
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						HTTPSOnly: util.Ref(false),
						SiteConfig: &armappservice.SiteConfig{
							MinTLSVersion: nil,
						},
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{},
					},
				},
			},
			want: &ontology.Function{
				Id:           "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.web/sites/webapp2",
				Name:         "WebApp2",
				CreationTime: nil,
				Labels: map[string]string{
					"testKey1": "testTag1",
					"testKey2": "testTag2",
				},
				GeoLocation: &ontology.GeoLocation{
					Region: "West Europe",
				},
				ParentId:            util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1"),
				Raw:                 "{\"*armappservice.Site\":[{\"id\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2\",\"kind\":\"app\",\"location\":\"West Europe\",\"name\":\"WebApp2\",\"properties\":{\"httpsOnly\":false,\"resourceGroup\":\"res1\",\"siteConfig\":{},\"virtualNetworkSubnetId\":\"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2\"},\"tags\":{\"testKey1\":\"testTag1\",\"testKey2\":\"testTag2\"}}],\"armappservice.WebAppsClientGetConfigurationResponse\":[{\"properties\":{}}]}",
				NetworkInterfaceIds: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/virtualnetworks/vnet1/subnets/subnet2"},
				ResourceLogging: &ontology.ResourceLogging{
					Enabled: false,
				},
				/*HttpEndpoint: &ontology.HttpEndpoint{
					TransportEncryption: &ontology.TransportEncryption{
						Enabled:    false,
						Enforced:   false,
						TlsVersion: "",
						Algorithm:  "",
					},
				},*/
				InternetAccessibleEndpoint: false,
				Redundancies:               []*ontology.Redundancy{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery

			// Set clients if needed
			if tt.fields.clientWebApps {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			assert.Equal(t, tt.want, d.handleWebApp(tt.args.webApp, tt.args.config))
		})
	}
}

func Test_getTransportEncryption(t *testing.T) {
	type args struct {
		siteProps *armappservice.SiteProperties
		config    armappservice.WebAppsClientGetConfigurationResponse
	}
	tests := []struct {
		name    string
		args    args
		wantEnc *ontology.TransportEncryption
	}{
		{
			name: "Happy path: TLSVersion/CipherSuite not available",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(false),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{},
					},
				},
			},
			wantEnc: &ontology.TransportEncryption{
				Enforced: false,
				Enabled:  false,
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.0, TLS version 1.0",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne0),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &ontology.TransportEncryption{
				Enforced:        true,
				Enabled:         true,
				Protocol:        constants.TLS,
				ProtocolVersion: 1.0,
				CipherSuites: []*ontology.CipherSuite{
					{
						SessionCipher: "AES-128-GCM",
						MacAlgorithm:  "SHA-256",
					},
				},
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.0, TLS version 1.1",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne1),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &ontology.TransportEncryption{
				Enforced:        true,
				Enabled:         true,
				Protocol:        constants.TLS,
				ProtocolVersion: 1.1,
				CipherSuites: []*ontology.CipherSuite{
					{
						SessionCipher: "AES-128-GCM",
						MacAlgorithm:  "SHA-256",
					},
				},
			},
		},
		{
			name: "Happy path: TLSVersion/CipherSuite available, TLS version 1.2",
			args: args{
				siteProps: &armappservice.SiteProperties{
					SiteConfig: &armappservice.SiteConfig{},
					HTTPSOnly:  util.Ref(true),
				},
				config: armappservice.WebAppsClientGetConfigurationResponse{
					SiteConfigResource: armappservice.SiteConfigResource{
						Properties: &armappservice.SiteConfig{
							MinTLSVersion:     util.Ref(armappservice.SupportedTLSVersionsOne2),
							MinTLSCipherSuite: util.Ref(armappservice.TLSCipherSuitesTLSAES128GCMSHA256),
						},
					},
				},
			},
			wantEnc: &ontology.TransportEncryption{
				Enforced:        true,
				Enabled:         true,
				Protocol:        constants.TLS,
				ProtocolVersion: 1.2,
				CipherSuites: []*ontology.CipherSuite{
					{
						SessionCipher: "AES-128-GCM",
						MacAlgorithm:  "SHA-256",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEnc := getTransportEncryption(tt.args.siteProps, tt.args.config)
			assert.Equal(t, tt.wantEnc, gotEnc)
		})
	}
}

func Test_azureComputeDiscovery_getResourceLoggingWebApp(t *testing.T) {
	type fields struct {
		azureDiscovery     *azureDiscovery
		defenderProperties map[string]*defenderProperties
		clientWebApp       bool
	}
	type args struct {
		site *armappservice.Site
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantRl *ontology.ResourceLogging
	}{
		{
			name: "Input empty",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
			},
			args: args{
				site: nil,
			},
			wantRl: &ontology.ResourceLogging{},
		},
		{
			name: "Happy path: logging disabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApp:   true,
			},
			args: args{
				site: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp2"),
					Name: util.Ref("WebApp2"),
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess:    util.Ref("Enabled"),
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
			},
			wantRl: &ontology.ResourceLogging{
				Enabled: false,
			},
		},
		{
			name: "Happy path: logging enabled",
			fields: fields{
				azureDiscovery: NewMockAzureDiscovery(newMockSender()),
				clientWebApp:   true,
			},
			args: args{
				site: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Kind: util.Ref("app"),
					Properties: &armappservice.SiteProperties{
						PublicNetworkAccess:    util.Ref("Enabled"),
						ResourceGroup:          util.Ref("res1"),
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet2"),
					},
				},
			},
			wantRl: &ontology.ResourceLogging{
				Enabled: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.azureDiscovery
			d.defenderProperties = tt.fields.defenderProperties

			// Set clients if needed
			if tt.fields.clientWebApp {
				// initialize webApps client
				_ = d.initWebAppsClient()
			}

			gotRl := d.getResourceLoggingWebApps(tt.args.site)
			assert.Equal(t, tt.wantRl, gotRl)
		})
	}
}

func Test_getRedundancies(t *testing.T) {
	type args struct {
		app *armappservice.Site
	}
	tests := []struct {
		name string
		args args
		want []*ontology.Redundancy
	}{
		{
			name: "Happy path: no redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeNone),
					},
				},
			},
			want: nil,
		},
		{
			name: "Happy path: zone redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeActiveActive),
					},
				},
			},
			want: []*ontology.Redundancy{
				{Type: &ontology.Redundancy_ZoneRedundancy{}},
			},
		},
		{
			name: "Happy path: zone and geo redundancy",
			args: args{
				app: &armappservice.Site{
					ID:   util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Web/sites/WebApp1"),
					Name: util.Ref("WebApp1"),
					Properties: &armappservice.SiteProperties{
						RedundancyMode: util.Ref(armappservice.RedundancyModeGeoRedundant),
					},
				},
			},
			want: []*ontology.Redundancy{
				{Type: &ontology.Redundancy_ZoneRedundancy{}},
				{Type: &ontology.Redundancy_GeoRedundancy{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getRedundancies(tt.args.app)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_publicNetworkAccessStatus(t *testing.T) {
	type args struct {
		status *string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty input",
			args: args{},
			want: false,
		},
		{
			name: "Happy path: Enabled",
			args: args{
				status: util.Ref("Enabled"),
			},
			want: true,
		},
		{
			name: "Happy path: Empty String",
			args: args{
				status: util.Ref(""),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := publicNetworkAccessStatus(tt.args.status); got != tt.want {
				t.Errorf("publicNetworkAccessStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVirtualNetworkSubnetId(t *testing.T) {
	type args struct {
		site *armappservice.Site
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty input",
			args: args{},
			want: []string{},
		},
		{
			name: "Happy path: with virtual network subnet ID",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{
						VirtualNetworkSubnetID: util.Ref("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"),
					},
				},
			},
			want: []string{"/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/res1/providers/microsoft.network/virtualnetworks/vnet1/subnets/subnet1"},
		},
		{
			name: "Happy path: without virtual network subnet ID",
			args: args{
				site: &armappservice.Site{
					Properties: &armappservice.SiteProperties{},
				},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVirtualNetworkSubnetId(tt.args.site)
			assert.Equal(t, tt.want, got)
		})
	}
}
