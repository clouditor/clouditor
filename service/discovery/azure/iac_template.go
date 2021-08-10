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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func (d *azureIacTemplateDiscovery) Name() string {
	return "Azure"
}

func (d *azureIacTemplateDiscovery) Description() string {
	return "Discovery IaC template."
}

// Discover IaC template
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

		template, ok := result.Template.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("IaC template type convertion failed")
		}

		err = saveExportTemplate(result, *resourceGroups[i].Name)
		if err != nil {
			fmt.Println("Error saving export template: ", err)
		}

		for templateKey, templateValue := range template {

			if templateKey == "resources" {
				resources, ok := templateValue.([]interface{})
				if !ok {
					return nil, fmt.Errorf("templateValue  type convertion failed")
				}

				for _, resourcesValue := range resources {
					value, ok := resourcesValue.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("resources type convertion failed")
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
							} else if valueValue.(string) == "Microsoft.Storage/storageAccounts" {
								storage, err := d.createStorageResource(value, *resourceGroups[i].Name)
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
func saveExportTemplate(template resources.GroupExportResult, groupName string) error {

	var (
		filepath     string
		filename     string
		fileTemplate string
	)

	prefix, indent := "", "    "
	exported, err := json.MarshalIndent(template, prefix, indent)
	if err != nil {
		return fmt.Errorf("MarshalIndent failed %w", err)
	}

	filepath = "../../results/raw_discovery_results/azure_iac_raw_templates/"
	fileTemplate = "%s-template.json"
	filename = fmt.Sprintf(fileTemplate, groupName)

	// Check if folder exists
	err = os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("check for directory existence failed:  %w", err)
	}

	err = ioutil.WriteFile(filepath+filename, exported, 0666)
	if err != nil {
		return fmt.Errorf("write file failed %w", err)
	} else {
		log.Infof("raw IaC template file written to: {%s}{%s}", filepath, filename)

	}

	return nil
}

func (d *azureIacTemplateDiscovery) createStorageResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {

	var (
		name string
	)

	resourceType := resourceValue["type"].(string)

	for key, value := range resourceValue {
		// Get storage account name
		if key == "name" {
			name = getDefaultNameOfResource(value.(string))
		}
	}

	storage := &voc.ObjectStorage{
		Storage: &voc.Storage{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(d.createID(resourceGroup, resourceType, name)),
				Name:         name,
				CreationTime: 0, // No creation time available
				Type:         []string{"ObjectStorage", "Storage", "Resource"},
			},
			AtRestEncryption: &voc.AtRestEncryption{
				Keymanager: getStorageKeySource(resourceValue),
				Algorithm:  "AES-265", // seems to be always AWS-256,
				Enabled:    blobServiceEncryptionEnabled(resourceValue),
			},
		},
		HttpEndpoint: &voc.HttpEndpoint{
			Url:           "", // Not able to get from IaC template
			Functionality: &voc.Functionality{},
			Authenticity: &voc.Authenticity{
				SecurityFeature: &voc.SecurityFeature{},
			},
			TransportEncryption: &voc.TransportEncryption{
				Enabled:    true, // cannot be disabled
				Enforced:   httpTrafficOnlyEnabled(resourceValue),
				TlsVersion: getMinTlsVersion(resourceValue),
				Algorithm:  "",
			},
			Method:  "",
			Handler: "",
			Path:    "",
		},
	}

	return storage, nil
}

func httpTrafficOnlyEnabled(value map[string]interface{}) bool {

	if httpTrafficOnlyEnabled, ok := value["properties"].(map[string]interface{})["supportsHttpsTrafficOnly"].(bool); ok {
		return httpTrafficOnlyEnabled
	}

	return false
}

func getStorageKeySource(value map[string]interface{}) string {

	if storageKeySource, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["keySource"].(string); ok {
		return storageKeySource
	}

	return ""
}

func blobServiceEncryptionEnabled(value map[string]interface{}) bool {

	if blobServiceEnabled, ok := value["properties"].(map[string]interface{})["encryption"].(map[string]interface{})["services"].(map[string]interface{})["blob"].(map[string]interface{})["enabled"].(bool); ok {
		return blobServiceEnabled
	}

	return false
}

func getMinTlsVersion(value map[string]interface{}) string {

	if minTlsVersion, ok := value["properties"].(map[string]interface{})["minimumTlsVersion"].(string); ok {
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
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false,
			RestrictedPorts: "",
		},
		HttpEndpoint: &[]voc.HttpEndpoint{},
		// // TODO(all): Do we need the httpEndpoint?
		// HttpEndpoint: &voc.HttpEndpoint{},
	}

	return lb, nil
}

func (d *azureIacTemplateDiscovery) createVMResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {
	var id string
	var name string
	var enabled bool

	for key, value := range resourceValue {

		// Get VM name
		if key == "name" {
			name = getDefaultNameOfResource(value.(string))
		}

		// Get bool for Logging enabled
		if key == "properties" {
			properties, ok := value.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("type convertion failed")
			}

			for propertiesKey, propertiesValue := range properties {
				if propertiesKey == "diagnosticsProfile" {
					enabled = propertiesValue.(map[string]interface{})["bootDiagnostics"].(map[string]interface{})["enabled"].(bool)
				}
			}
		}
	}

	// Get ID
	// ID must be put together by hand, is not available in template. Better ideas? Leave empty?
	id = d.createID(resourceGroup, resourceValue["type"].(string), name)

	vm := &voc.VirtualMachine{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID:           voc.ResourceID(id),
				Name:         name,
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
			}},
		Log: &voc.Log{
			Activated: enabled,
		},
	}

	return vm, nil
}

func (d *azureIacTemplateDiscovery) createID(resourceGroup, resourceType, name string) string {
	return "/subscriptions/" + *d.sub.SubscriptionID + "/resourceGroups/" + strings.ToUpper(resourceGroup) + "/providers/" + resourceType + "/" + name
}

// getDefaultNameOfResource gets the defaultName from template parameter
// TODO(all): The exported template contains a parameter instead of the defaultName (resourceName). Furthermore, the template parameters do not contain a mapping from the parameter to the defaultName. In the parameter name all word separators (e.g. _, -, .) were replaced by a underscore (_), so it is not possible to uniquely restore the defaultName. Ideas? Do we need the correct defaultNames?
func getDefaultNameOfResource(name string) string {
	// Name in template is an parameter and unnecessary information must be shortened
	nameSplit := strings.Split(name, "'")
	anotherNameSplit := strings.Split(nameSplit[1], "_")
	anotherNameSplit = anotherNameSplit[1:]
	anotherNameSplit = anotherNameSplit[:len(anotherNameSplit)-1]
	resourceDefaultName := strings.Join(anotherNameSplit, "-")

	return resourceDefaultName
}
