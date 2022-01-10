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
	"context"
	"errors"
	"fmt"
	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-02-01/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureARMTemplateDiscovery struct {
	azureDiscovery
}

func NewAzureARMTemplateDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureARMTemplateDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (*azureARMTemplateDiscovery) Name() string {
	return "Azure ARM template"
}

func (*azureARMTemplateDiscovery) Description() string {
	return "Discovery using an Azure Resource Manager (ARM) template."
}

// List Azure resources by discovering Azure ARM template
func (d *azureARMTemplateDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	armResources, err := d.discoverARMTemplate()
	if err != nil {
		return nil, fmt.Errorf("could not discover Azure ARM template: %w", err)
	}
	list = append(list, armResources...)

	return
}

func (d *azureARMTemplateDiscovery) discoverARMTemplate() ([]voc.IsCloudResource, error) {

	var (
		list []voc.IsCloudResource
	)

	client := resources.NewGroupsClient(to.String(d.sub.SubscriptionID))
	d.apply(&client.Client)

	// Get all resource groups, as the exportTemplate function works only on resource group level
	resultResourceGroupsResult, err := client.ListComplete(context.Background(), "", nil)
	if err != nil {
		return nil, fmt.Errorf("could not discover resource groups: %w", err)
	}

	resourceGroups := *resultResourceGroupsResult.Response().Value
	for i := range resourceGroups {

		expReq := resources.ExportTemplateRequest{
			ResourcesProperty: &[]string{"*"},
		}
		result, err := client.ExportTemplate(context.Background(), *resourceGroups[i].Name, expReq)
		if err != nil {
			return nil, fmt.Errorf("could not discover Azure ARM template: %w", err)
		}

		// TODO Update to latest API version -> GroupsExportTemplateFuture
		// result = GroupExportResult
		armTemplate, ok := result.Template.(map[string]interface{})
		if !ok {
			return nil, errors.New("ARM template type assertion failed")
		}

		for templateKey, templateValue := range armTemplate {
			if templateKey == "resources" {
				azureResource, ok := templateValue.([]interface{})
				if !ok {
					return nil, errors.New("templateValue type assertion failed")
				}

				for _, resourcesValue := range azureResource {
					value, ok := resourcesValue.(map[string]interface{})
					if !ok {
						return nil, errors.New("azureResource type assertion failed")
					}

					for valueKey, valueValue := range value {
						if valueKey == "type" {

							if valueValue.(string) == "Microsoft.Compute/virtualMachines" {
								vm, err := d.handleVirtualMachine(armTemplate, value, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create virtual machine resource: %w", err)
								}
								list = append(list, vm)
							} else if valueValue.(string) == "Microsoft.Network/loadBalancers" {
								lb, err := d.handleLoadBalancer(armTemplate, value, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create load balancer resource: %w", err)
								}
								list = append(list, lb)
							} else if valueValue.(string) == "Microsoft.Storage/storageAccounts/blobServices/containers" {
								storage, err := d.handleObjectStorage(value, azureResource, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create storage resource: %w", err)
								}
								list = append(list, storage)
							} else if valueValue.(string) == "Microsoft.Storage/storageAccounts/fileServices/shares" {
								storage, err := d.handleFileStorage(value, azureResource, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create storage resource: %w", err)
								}
								list = append(list, storage)
							}
							// TODO(garuppel): Handle BlockStorage resources?
						}
					}
				}
			}
		}
	}

	return list, nil
}

func (d *azureARMTemplateDiscovery) handleObjectStorage(resourceValue map[string]interface{}, azureResources []interface{}, resourceGroup string) (voc.IsCompute, error) {

	var (
		azureResourceName string
		storage           voc.IsCompute
		enc               voc.HasAtRestEncryption
	)

	// The resources are only referencing to parameters instead of using the resource names
	// In case of object storages, we take the container name as resource name
	azureResourceName = containerNameFromResourceName(resourceValue["name"].(string))

	// 'dependsOn' references to the related Azure ARM template resources. For the storage account information, we need the related storage account resource name
	dependsOnList, ok := (resourceValue["dependsOn"]).([]interface{})
	if !ok {
		return nil, errors.New("dependsOn type assertion failed")
	}

	storageAccountResource, err := storageAccountResourceFromARMTemplate(dependsOnList, azureResources)
	if err != nil {
		return nil, fmt.Errorf("cannot get storage account resource from Azure ARM template: %w", err)
	}

	enc, err = storageAccountAtRestEncryptionFromARMtemplate(storageAccountResource)
	if err != nil {
		return nil, fmt.Errorf("cannot get atRestEncryption for storage account resource from Azure ARM template: #{err)")
	}

	storage = &voc.ObjectStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(d.createID(resourceGroup, "blobServices", azureResourceName)),
				Name:         azureResourceName,
				CreationTime: 0, // No creation time available
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: storageAccountResource["location"].(string),
				},
			},
			AtRestEncryption: enc,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url: "", // not available
			TransportEncryption: &voc.TransportEncryption{
				Enabled:    isServiceEncryptionEnabled("blob", storageAccountResource),
				Enforced:   isHttpsTrafficOnlyEnabled(storageAccountResource),
				TlsVersion: minTlsVersionOfStorageAccount(storageAccountResource),
				Algorithm:  "TLS",
			},
		},
	}

	return storage, nil
}

func (d *azureARMTemplateDiscovery) handleFileStorage(resourceValue map[string]interface{}, azureResources []interface{}, resourceGroup string) (voc.IsCompute, error) {

	var (
		azureResourceName string
		storage           voc.IsCompute
		enc               voc.HasAtRestEncryption
	)

	// The resources are only referencing to parameters instead of using the resource names
	azureResourceName = containerNameFromResourceName(resourceValue["name"].(string))

	// Necessary to get the needed information from the Azure ARM template
	dependsOnList, ok := (resourceValue["dependsOn"]).([]interface{})
	if !ok {
		return nil, errors.New("dependsOn type assertion failed")
	}

	storageAccountResource, err := storageAccountResourceFromARMTemplate(dependsOnList, azureResources)
	if err != nil {
		return nil, fmt.Errorf("cannot get storage account resource from Azure ARM template: %w", err)
	}

	enc, err = storageAccountAtRestEncryptionFromARMtemplate(storageAccountResource)
	if err != nil {
		return nil, fmt.Errorf("cannot get atRestEncryption for storage account resource from Azure ARM template: %w", err)
	}

	storage = &voc.FileStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(d.createID(resourceGroup, "fileServices", azureResourceName)),
				Name:         azureResourceName,
				CreationTime: 0, // No creation time available
				Type:         []string{"FileStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: storageAccountResource["location"].(string),
				},
			},
			AtRestEncryption: enc,
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url: "", // not available
			TransportEncryption: &voc.TransportEncryption{
				Enabled:    isServiceEncryptionEnabled("file", storageAccountResource),
				Enforced:   isHttpsTrafficOnlyEnabled(storageAccountResource),
				TlsVersion: minTlsVersionOfStorageAccount(storageAccountResource),
				Algorithm:  "TLS",
			},
		},
	}

	return storage, nil
}

func storageAccountAtRestEncryptionFromARMtemplate(storageAccountResource map[string]interface{}) (voc.HasAtRestEncryption, error) {

	var enc voc.HasAtRestEncryption

	encType, ok := storageAccountResource["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["keySource"].(string)

	if !ok {
		return nil, errors.New("type assertion failed")
	}

	if encType == "Microsoft.Storage" {
		enc = voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "AES256",
				Enabled:   isServiceEncryptionEnabled("blob", storageAccountResource),
			},
		}
	} else if encType == "Microsoft.Keyvault" {
		keyVaultUrl := storageAccountResource["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["keyvaultproperties"].(map[string]interface{})["keyvaulturi"].(string)

		enc = voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "", // TODO(garuppel): TBD
				Enabled:   isServiceEncryptionEnabled("blob", storageAccountResource),
			},
			KeyUrl: keyVaultUrl,
		}
	}

	return enc, nil
}

func storageAccountResourceFromARMTemplate(resourceNames []interface{}, azureTemplateResources []interface{}) (map[string]interface{}, error) {

	var (
		resourceType string
		resourceName string
	)

	// Get parameter type and name from corresponding storage account resource
	for _, resourceNameElem := range resourceNames {

		resourceNameSplit := strings.Split(resourceNameElem.(string), ",")
		if len(resourceNameSplit) == 2 {
			resourceType = strings.Split(resourceNameSplit[0], "'")[1]
			resourceName = "[" + resourceNameSplit[1][1:len(resourceNameSplit[1])-2] + "]"
		}
	}

	for _, resourcesValue := range azureTemplateResources {
		templateResources, ok := resourcesValue.(map[string]interface{})
		if !ok {
			return nil, errors.New("type assertion failed")
		}

		if templateResources["type"] == resourceType && templateResources["name"] == resourceName {
			return resourcesValue.(map[string]interface{}), nil
		}
	}

	return nil, errors.New("could not get resource from Azure ARM template")
}

func isHttpsTrafficOnlyEnabled(value map[string]interface{}) bool {

	if supportsHttpsTrafficOnly, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["supportsHttpsTrafficOnly"].(bool); ok {
		return supportsHttpsTrafficOnly
	}

	return false
}

func isServiceEncryptionEnabled(serviceType string, value map[string]interface{}) bool {

	if blobServiceEnabled, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["services"].(map[string]interface{})[serviceType].(map[string]interface{})["enabled"].(bool); ok {
		return blobServiceEnabled
	}

	return false
}

func minTlsVersionOfStorageAccount(value map[string]interface{}) string {

	if minTlsVersion, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["minimumTlsVersion"].(string); ok {
		return minTlsVersion
	}

	return ""
}

func (d *azureARMTemplateDiscovery) handleLoadBalancer(template map[string]interface{}, resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {

	var name string
	var err error

	resourceType := resourceValue["type"].(string)

	for key, value := range resourceValue {
		// Get LB name
		if key == "name" {
			name, err = defaultResourceNameFromParameter(template, value.(string))
			if err != nil {
				return nil, errors.New("getting parameter default name failed")
			}
		}
	}

	// TODO(garuppel): Which additional information do we get from the template?
	lb := &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				CloudResource: &voc.CloudResource{
					ID:           voc.ResourceID(d.createID(resourceGroup, resourceType, name)),
					Name:         name,
					CreationTime: 0, // No creation time available
					Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
					GeoLocation: voc.GeoLocation{
						Region: resourceValue["location"].(string),
					},
				},
			},
			Compute: []voc.ResourceID{},
			Ips:     []string{},
			Ports:   []int16{},
		},
		AccessRestrictions: &[]voc.AccessRestriction{
			// TODO(all)
			//Inbound:         false,
			//RestrictedPorts: "",
		},
		// TODO(all): Do we need the httpEndpoint?
		HttpEndpoints:   &[]voc.HttpEndpoint{},
		Url:             "",                 // TODO(all): TBD
		NetworkServices: []voc.ResourceID{}, // TODO(all): TBD
	}

	return lb, nil
}

func (d *azureARMTemplateDiscovery) handleVirtualMachine(template map[string]interface{}, resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {
	var id string
	var name string
	var bootDiagnosticsEnabled bool
	var storageUri string
	var properties map[string]interface{}
	var ok bool
	var err error

	for key, value := range resourceValue {

		// Get VM name
		if key == "name" {
			name, err = defaultResourceNameFromParameter(template, value.(string))
			if err != nil {
				return nil, errors.New("getting parameter default name failed")
			}
		}

		// Get boot logging status (bootDiagnosticsEnabled)
		if key == "properties" {
			properties, ok = value.(map[string]interface{})

			if !ok {
				return nil, errors.New("type assertion failed")
			}

			for propertiesKey, propertiesValue := range properties {
				if propertiesKey == "diagnosticsProfile" {
					bootDiagnosticsEnabled = propertiesValue.(map[string]interface{})["bootDiagnostics"].(map[string]interface{})["enabled"].(bool)
					storageUri = storageURIFromARMTemplate(propertiesValue.(map[string]interface{}))
				}
			}
		}
	}

	// Get virtual machine ID
	// Virtual machine ID must be put together by hand, is not available in template
	id = d.createID(resourceGroup, resourceValue["type"].(string), name)

	vm := &voc.VirtualMachine{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(id),
				Name:         name,
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: resourceValue["location"].(string),
				},
			},
		},
		BootLog: &voc.BootLog{
			Log: &voc.Log{
				Enabled:         bootDiagnosticsEnabled,
				Output:          []voc.ResourceID{voc.ResourceID(storageUri)},
				RetentionPeriod: 0, // Currently, configuring the retention period for Managed Boot Diagnostics is not available. The logs will be overwritten after 1gb of space according to https://github.com/MicrosoftDocs/azure-docs/issues/69953
			},
		},
		OSLog: &voc.OSLog{
			Log: &voc.Log{
				Enabled: false,
			}, // TODO(all): available in ARM template/Azure?
		},
		BlockStorage: d.blockStorageResourceIDs(properties, resourceGroup),
	}

	return vm, nil
}

func (d *azureARMTemplateDiscovery) blockStorageResourceIDs(properties map[string]interface{}, resourceGroupName string) []voc.ResourceID {
	var blockStorage []voc.ResourceID

	dataDisks := properties["storageProfile"].(map[string]interface{})["dataDisks"].([]interface{})
	for _, dataDisk := range dataDisks {
		dataDiskName := dataDisk.(map[string]interface{})["name"].(string)
		dataDiskResourceId := d.createID(resourceGroupName, "Microsoft.Compute/disks", dataDiskName)
		blockStorage = append(blockStorage, voc.ResourceID(dataDiskResourceId))
	}

	return blockStorage
}

func (d *azureARMTemplateDiscovery) createID(resourceGroup, resourceType, name string) string {
	return "/subscriptions/" + *d.sub.SubscriptionID + "/resourceGroups/" + resourceGroup + "/providers/" + resourceType + "/" + name
}

// defaultResourceNameFromParameter returns the default name given in the parameters section of the template. If not possible, get name from parameter name, e.g., [parameters('loadBalancers_kubernetes_name')]
func defaultResourceNameFromParameter(template map[string]interface{}, name string) (string, error) {

	// [parameters('loadBalancers_kubernetes_name')]
	for templateKey, templateValue := range template {
		if templateKey == "parameters" {
			resourceParameter, ok := templateValue.(map[string]interface{})
			if !ok {
				return "", errors.New("templateValue type assertion failed")
			}

			resource, ok := resourceParameter[coreNameFromResourceName(name)].(map[string]interface{})
			if !ok {
				return "", errors.New("parameter resource type assertion failed")
			}

			if resource["defaultValue"] == nil {
				return coreNameFromResourceName(name), nil
			}

			return resource["defaultValue"].(string), nil
		}
	}

	return "", errors.New("error getting default resource name")
}

// coreNameFromResourceName returns the parameter name without the additional information around. Necessary if in parameters no default name exists.
// Example: [parameters('virtualMachines_vm3_name')] returns 'vm3'
func coreNameFromResourceName(name string) string {
	return strings.Split(name, "'")[1]
}

// containerNameFromResourceName returns the container name of the resource type 'Microsoft.Storage/storageAccounts/blobServices/containers'
// Example: [concat(parameters('storageAccounts_storage1_name'), 'default/container1')] returns 'container1'
func containerNameFromResourceName(name string) string {
	nameSplit := strings.Split(name, "'")
	anotherNameSplit := strings.Split(nameSplit[3], "/")

	return anotherNameSplit[1]
}

func storageURIFromARMTemplate(name map[string]interface{}) string {
	var storageUri string
	if name["bootDiagnostics"].(map[string]interface{})["storageUri"] == nil {
		return ""
	}

	// URI example: "[concat('https://', parameters('storageAccounts_test_name'), '.blob.core.windows.net/')]"
	storageUriInformation := name["bootDiagnostics"].(map[string]interface{})["storageUri"].(string)
	nameSplit := strings.Split(storageUriInformation, "'")
	storageUriName := strings.Split(nameSplit[3], "_")
	storageUri = nameSplit[1] + storageUriName[1] + nameSplit[5]

	return storageUri
}
