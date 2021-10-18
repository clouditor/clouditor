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
	"fmt"

	"strings"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-02-01/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

type azureIacTemplateDiscovery struct {
	azureDiscovery
}

func NewAzureIacTemplateDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureIacTemplateDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (*azureIacTemplateDiscovery) Name() string {
	return "Azure"
}

func (*azureIacTemplateDiscovery) Description() string {
	return "Discovery IaC template."
}

// List Azure resources by discoverying IaC template
func (d *azureIacTemplateDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("could not authorize Azure account: %w", err)
	}

	iacTemplate, err := d.discoverIaCTemplate()
	if err != nil {
		return nil, fmt.Errorf("could not discover IaC template: %w", err)
	}
	list = append(list, iacTemplate...)

	return
}

func (d *azureIacTemplateDiscovery) discoverIaCTemplate() ([]voc.IsCloudResource, error) {

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
			return nil, fmt.Errorf("could not discover IaC template: %w", err)
		}

		//err = saveExportTemplate(result, *resourceGroups[i].Name)
		//if err != nil {
		//	return nil, fmt.Errorf("could not save IaC template: %w", err)
		//}

		template, ok := result.Template.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("IaC template type convertion failed")
		}

		for templateKey, templateValue := range template {
			if templateKey == "resources" {
				azureResource, ok := templateValue.([]interface{})
				if !ok {
					return nil, fmt.Errorf("templateValue type convertion failed")
				}

				for _, resourcesValue := range azureResource {
					value, ok := resourcesValue.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("azureResource type convertion failed")
					}

					for valueKey, valueValue := range value {
						if valueKey == "type" {

							if valueValue.(string) == "Microsoft.Compute/virtualMachines" {
								vm, err := d.createVMResource(value, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create virtual machine resource: %w", err)
								}
								list = append(list, vm)
							} else if valueValue.(string) == "Microsoft.Network/loadBalancers" {
								lb, err := d.createLBResource(value, *resourceGroups[i].Name)
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
							} else if valueValue.(string) == "Microsoft.Storage/storageAccounts/fileServices/shares" {
								storage, err := d.handleFileStorage(value, azureResource, *resourceGroups[i].Name)
								if err != nil {
									return nil, fmt.Errorf("could not create storage resource: %w", err)
								}
								list = append(list, storage)
							}
						}
					}
				}
			}
		}
	}

	return list, nil
}

// saveExportTemplate saves the resource group template in a json file.
//func saveExportTemplate(template resources.GroupExportResult, groupName string) error {
//
//	var (
//		filepath     string
//		filename     string
//		fileTemplate string
//	)
//
//	prefix, indent := "", "    "
//	exported, err := json.MarshalIndent(template, prefix, indent)
//	if err != nil {
//		return fmt.Errorf("MarshalIndent failed %w", err)
//	}
//
//	filepath = "../../results/raw_discovery_results/azure_iac_raw_templates/"
//	fileTemplate = "%s-template.json"
//	filename = fmt.Sprintf(fileTemplate, groupName)
//
//	// Check if folder exists
//	err = os.MkdirAll(filepath, os.ModePerm)
//	if err != nil {
//		return fmt.Errorf("check for directory existence failed:  %w", err)
//	}
//
//	err = ioutil.WriteFile(filepath+filename, exported, 0666)
//	if err != nil {
//		return fmt.Errorf("write file failed %w", err)
//	} else {
//		log.Infof("raw IaC template file written to: {%s}{%s}", filepath, filename)
//
//	}
//
//	return nil
//}

// TODO(garuppel): split to different storage types
func (d *azureIacTemplateDiscovery) handleObjectStorage(resourceValue map[string]interface{}, azureResources []interface{}, resourceGroup string) (voc.IsCompute, error) {

	var (
		azureResourceName string
		storage voc.IsCompute
		enc voc.HasAtRestEncryption
	)

	// The resources are only referencing to parameters instead of using the resource names
	azureResourceName = getDefaultNameOfStorageResource(resourceValue["name"].(string))

	// Necessary to get the needed information from the IaC template
	dependsOnList, ok :=(resourceValue["dependsOn"]).([]interface{})
	if !ok {
		return nil, fmt.Errorf("dependsOn  convertion failed")
	}

	storageAccountResource, err := getStorageAccountResourceFromTemplate(dependsOnList, azureResources)
	if err != nil {
		return nil, fmt.Errorf("cannot create object storage")
	}

	enc = getObjectStorageAtRestEncryptionFromIac(storageAccountResource)

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
				Enabled:    true, // TODO get from IaC template
				Enforced:   httpsTrafficOnlyEnabled(storageAccountResource),
				TlsVersion: getMinTlsVersionOfStorageAccount(storageAccountResource),
				Algorithm:  "", // not available
			},
		},
	}

	return storage, nil
}

func (d *azureIacTemplateDiscovery) handleFileStorage(resourceValue map[string]interface{}, azureResources []interface{}, resourceGroup string) (voc.IsCompute, error) {

	var (
		azureResourceName string
		storage voc.IsCompute
		enc voc.HasAtRestEncryption
	)

	// The resources are only referencing to parameters instead of using the resource names
	azureResourceName = getDefaultNameOfStorageResource(resourceValue["name"].(string))

	// Necessary to get the needed information from the IaC template
	dependsOnList, ok :=(resourceValue["dependsOn"]).([]interface{})
	if !ok {
		return nil, fmt.Errorf("dependsOn  convertion failed")
	}

	storageAccountResource, err := getStorageAccountResourceFromTemplate(dependsOnList, azureResources)
	if err != nil {
		return nil, fmt.Errorf("cannot create object storage")
	}

	enc = getObjectStorageAtRestEncryptionFromIac(storageAccountResource)

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
				Enabled:    true, // TODO get from IaC template
				Enforced:   httpsTrafficOnlyEnabled(storageAccountResource),
				TlsVersion: getMinTlsVersionOfStorageAccount(storageAccountResource),
				Algorithm:  "", // not available
			},
		},
	}

	return storage, nil
}


func getObjectStorageAtRestEncryptionFromIac(storageAccountResource map[string]interface{}) voc.HasAtRestEncryption {

	var enc voc.HasAtRestEncryption

	encType := storageAccountResource["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["keySource"].(string)

	if encType == "Microsoft.Storage" {
		enc = &voc.ManagedKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "AES-265", // seems to be always AWS-256,
				Enabled: blobServiceEncryptionEnabled(storageAccountResource),
			},
		}
	} else if encType == "Microsoft.Keyvault"{
		keyVaultUrl := storageAccountResource["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["keyvaultProperties"].(map[string]interface{})["keyvaulturi"].(string)

		enc = &voc.CustomerKeyEncryption{
			AtRestEncryption: &voc.AtRestEncryption{
				Algorithm: "AES-265", // seems to be always AWS-256,
				Enabled: blobServiceEncryptionEnabled(storageAccountResource),
			},
			KeyUrl: keyVaultUrl,
		}
	}

	return enc
}

func getStorageAccountResourceFromTemplate(resourceNames []interface{}, azureTemplateResources  []interface{}) (map[string]interface{}, error) {

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
			return nil, fmt.Errorf("azureResource type convertion failed")
		}

		if templateResources["type"] == resourceType && templateResources["name"] == resourceName {
			return resourcesValue.(map[string]interface{}), nil
		}
	}

	return nil, fmt.Errorf("could not get resource from IaC template")
}

func httpsTrafficOnlyEnabled(value map[string]interface{}) bool {

	if supportsHttpsTrafficOnly, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["supportsHttpsTrafficOnly"].(bool); ok {
		return supportsHttpsTrafficOnly
	}

	return false
}

func blobServiceEncryptionEnabled(value map[string]interface{}) bool {

	if blobServiceEnabled, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["services"].(map[string]interface{})["blob"].(map[string]interface{})["enabled"].(bool); ok {
		return blobServiceEnabled
	}

	return false
}

func getMinTlsVersionOfStorageAccount(value map[string]interface{}) string {

	if minTlsVersion, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["minimumTlsVersion"].(string); ok {
		return minTlsVersion
	}

	return ""
}

func (d *azureIacTemplateDiscovery) createLBResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {

	var name string

	resourceType := resourceValue["type"].(string)

	for key, value := range resourceValue {
		// Get LB name
		if key == "name" {
			name = getDefaultNameOfResource(value.(string))
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
				},
			},
			Compute: []voc.ResourceID{},
			Ips:     []string{},
			Ports:   []int16{},
		},
		AccessRestrictions: &[]voc.AccessRestriction{
			//Inbound:         false,
			//RestrictedPorts: "",
		},
		// // TODO(all): Do we need the httpEndpoint?
		HttpEndpoints: &[]voc.HttpEndpoint{},
		Url:           "", // TODO(all): TBD
		//NetworkServices: , // TODO(all): TBD
	}

	return lb, nil
}

func (d *azureIacTemplateDiscovery) createVMResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {
	var id string
	var name string
	var bootDiagnosticsEnabled bool
	var properties map[string]interface{}
	var ok bool

	for key, value := range resourceValue {

		// Get VM name
		if key == "name" {
			name = getDefaultNameOfResource(value.(string))
		}

		// Get boot logging status (bootDiagnosticsEnabled)
		if key == "properties" {
			properties, ok = value.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("type convertion failed")
			}

			for propertiesKey, propertiesValue := range properties {
				if propertiesKey == "diagnosticsProfile" {
					bootDiagnosticsEnabled = propertiesValue.(map[string]interface{})["bootDiagnostics"].(map[string]interface{})["enabled"].(bool)
				}
			}
		}
	}

	// Get virtual machine ID
	// TODO(all): ID must be put together by hand, is not available in template. Better ideas? Leave empty?
	id = d.createID(resourceGroup, resourceValue["type"].(string), name)

	vm := &voc.VirtualMachine{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(id),
				Name:         name,
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
			},
		},
		BootLog: &voc.BootLog{
			Log: &voc.Log{
				Enabled: bootDiagnosticsEnabled,
			},
		},
		OSLog: &voc.OSLog{
			Log: &voc.Log{},
		},
		BlockStorage: d.getBlockStorageResourceIDs(properties, resourceGroup), 
	}

	return vm, nil
}

func (d *azureIacTemplateDiscovery) getBlockStorageResourceIDs(properties map[string]interface{}, resourceGroupName string) []voc.ResourceID {
	var blockStorage []voc.ResourceID

	dataDisks := properties["storageProfile"].(map[string]interface{})["dataDisks"].([]interface{})
	for _, dataDisk := range dataDisks {
		dataDiskName := dataDisk.(map[string]interface{})["name"].(string)
		dataDiskResourceId := d.createID(resourceGroupName, "Microsoft.Compute/disks", dataDiskName)
		blockStorage = append(blockStorage, voc.ResourceID(dataDiskResourceId))
	}

	return blockStorage
}

func (d *azureIacTemplateDiscovery) createID(resourceGroup, resourceType, name string) string {
	return "/subscriptions/" + *d.sub.SubscriptionID + "/resourceGroups/" + resourceGroup + "/providers/" + resourceType + "/" + name
}

// getDefaultNameOfResource gets the defaultName from template parameter
// TODO(all): The exported template contains a parameter instead of the defaultName (resourceName). Furthermore, the template parameters do not contain a mapping from the parameter to the defaultName. In the parameter name all word separators (e.g. _, -, .) were replaced by a underscore (_), so it is not possible to uniquely restore the defaultName. Ideas? Do we need the correct defaultNames?
func getDefaultNameOfResource(name string) string {
	// Name in template is a parameter and unnecessary information must be shortened
	nameSplit := strings.Split(name, "'")
	anotherNameSplit := strings.Split(nameSplit[1], "_")
	anotherNameSplit = anotherNameSplit[1:]
	anotherNameSplit = anotherNameSplit[:len(anotherNameSplit)-1]
	resourceDefaultName := strings.Join(anotherNameSplit, "-")

	return resourceDefaultName
}

func getDefaultNameOfStorageResource(name string) string {
	// Name in template is a parameter and unnecessary information must be shortened
	nameSplit := strings.Split(name, "'")
	anotherNameSplit := strings.Split(nameSplit[3], "/")

	return anotherNameSplit[1]
}