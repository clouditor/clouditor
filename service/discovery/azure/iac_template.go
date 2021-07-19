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
func (d *azureIacTemplateDiscovery) List() (list []voc.IsResource, err error) {
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

func (d *azureIacTemplateDiscovery) discoverIaCTemplate() ([]voc.IsResource, error) {

	var (
		list []voc.IsResource
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
			fmt.Println("ERROR: ", err)
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

	prefix, indent := "", "    "
	exported, err := json.MarshalIndent(template, prefix, indent)
	if err != nil {
		return fmt.Errorf("MarshalIndent failed %w", err)
	}

	fileTemplate := "%s-template.json"
	fileName := fmt.Sprintf(fileTemplate, groupName)
	if _, err := os.Stat(fileName); err == nil {
		return fmt.Errorf("file already exists")
	}

	ioutil.WriteFile(fileName, exported, 0666)
	if err != nil {
		return fmt.Errorf("write file failed %w", err)
	}

	return nil
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
	lb := &voc.LoadBalancerResource{
		NetworkService: voc.NetworkService{
			NetworkResource: voc.NetworkResource{
				Resource: voc.Resource{
					ID:           d.createID(resourceGroup, resourceType, name),
					Name:         name,
					CreationTime: 0, // No creation time available
					Type:         []string{"LoadBalancer", "NetworkService", "Resource"},
				},
			},
			IPs:   []string{},
			Ports: nil,
		},
		AccessRestriction: &voc.AccessRestriction{
			Inbound:         false,
			RestrictedPorts: "",
		},
		// TODO(all): Do we need the httpEndpoint?
		HttpEndpoints: []*voc.HttpEndpoint{},
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

	vm := &voc.VirtualMachineResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           id,
				Name:         name,
				CreationTime: 0, // No creation time available
				Type:         []string{"VirtualMachine", "Compute", "Resource"},
			}},
		Log: &voc.Log{
			Enabled: enabled,
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
