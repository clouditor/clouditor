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

type azureIaCTemplateDiscovery struct {
	azureDiscovery
}

func NewIaCTemplateDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureIaCTemplateDiscovery{}

	for _, opt := range opts {
		if auth, ok := opt.(*authorizerOption); ok {
			d.authOption = auth
		} else {
			d.options = append(d.options, opt)
		}
	}

	return d
}

func (d *azureIaCTemplateDiscovery) Name() string {
	return "Azure"
}

func (d *azureIaCTemplateDiscovery) Description() string {
	return "Discovery IaC template."
}

// Discover IaC template
func (d *azureIaCTemplateDiscovery) List() (list []voc.IsResource, err error) {
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

func (d *azureIaCTemplateDiscovery) discoverIaCTemplate() ([]voc.IsResource, error) {

	var list []voc.IsResource

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
			return nil, fmt.Errorf("could not discover IaC templates: %w", err)
		}

		template, ok := result.Template.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("type convertion failed")
		}

		for templateKey, templateValue := range template {
			if templateKey == "resources" {
				resources, ok := templateValue.([]interface{})
				if !ok {
					return nil, fmt.Errorf("type convertion failed")
				}
				for _, resourcesValue := range resources {
					value, ok := resourcesValue.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("type convertion failed")
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

func (d *azureIaCTemplateDiscovery) createLBResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {

	var name string

	resourceType := resourceValue["type"].(string)

	for key, value := range resourceValue {
		// Get LB name
		if key == "name" {
			name = getResourceName(value.(string))
		}
	}

	// TODO Which additional information do we get from the template?
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
		// TODO Do we need the httpEndpoint?
		HttpEndpoints: []*voc.HttpEndpoint{},
	}

	return lb, nil
}

func (d *azureIaCTemplateDiscovery) createVMResource(resourceValue map[string]interface{}, resourceGroup string) (voc.IsCompute, error) {
	var id string
	var name string
	var enabled bool

	for key, value := range resourceValue {

		// Get VM name
		if key == "name" {
			name = getResourceName(value.(string))
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

func (d *azureIaCTemplateDiscovery) createID(resourceGroup, resourceType, name string) string {
	return "/subscriptions/" + *d.sub.SubscriptionID + "/resourceGroups/" + strings.ToUpper(resourceGroup) + "/providers/" + resourceType + "/" + name
}

func getResourceName(name string) string {
	// Name in template is an parameter and unnecessary information must be shortened
	nameSplit := strings.Split(name, "'")
	vmNameSplit := strings.Split(nameSplit[1], "_")
	vmNameSplit = vmNameSplit[1:]
	vmNameSplit = vmNameSplit[:len(vmNameSplit)-1]
	// TODO Is it possible that an vm_name has a _ as delimiter
	resourceName := strings.Join(vmNameSplit, "-")

	return resourceName
}
