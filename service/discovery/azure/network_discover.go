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
	"clouditor.io/clouditor/v2/api/ontology"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

// discoverNetworkInterfaces discovers network interfaces
func (d *azureDiscovery) discoverNetworkInterfaces() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize network interfaces client
	if err := d.initNetworkInterfacesClient(); err != nil {
		return nil, err
	}

	// List all network interfaces
	err := listPager(d,
		d.clients.networkInterfacesClient.NewListAllPager,
		d.clients.networkInterfacesClient.NewListPager,
		func(res armnetwork.InterfacesClientListAllResponse) []*armnetwork.Interface {
			return res.Value
		},
		func(res armnetwork.InterfacesClientListResponse) []*armnetwork.Interface {
			return res.Value
		},
		func(ni *armnetwork.Interface) error {
			s := d.handleNetworkInterfaces(ni)

			log.Infof("Adding network interface '%s'", s.GetName())

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// discoverApplicationGateway discovers application gateways
func (d *azureDiscovery) discoverApplicationGateway() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize application gateway client
	if err := d.initApplicationGatewayClient(); err != nil {
		return nil, err
	}

	// List all application gateways
	err := listPager(d,
		d.clients.applicationGatewayClient.NewListAllPager,
		d.clients.applicationGatewayClient.NewListPager,
		func(res armnetwork.ApplicationGatewaysClientListAllResponse) []*armnetwork.ApplicationGateway {
			return res.Value
		},
		func(res armnetwork.ApplicationGatewaysClientListResponse) []*armnetwork.ApplicationGateway {
			return res.Value
		},
		func(ags *armnetwork.ApplicationGateway) error {
			s := d.handleApplicationGateway(ags)

			log.Infof("Adding application gateway %+v", s.GetName())

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}

// discoverLoadBalancer discovers load balancer
func (d *azureDiscovery) discoverLoadBalancer() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// initialize load balancers client
	if err := d.initLoadBalancersClient(); err != nil {
		return nil, err
	}

	// List all load balancers
	err := listPager(d,
		d.clients.loadBalancerClient.NewListAllPager,
		d.clients.loadBalancerClient.NewListPager,
		func(res armnetwork.LoadBalancersClientListAllResponse) []*armnetwork.LoadBalancer {
			return res.Value
		},
		func(res armnetwork.LoadBalancersClientListResponse) []*armnetwork.LoadBalancer {
			return res.Value
		},
		func(lbs *armnetwork.LoadBalancer) error {
			s := d.handleLoadBalancer(lbs)

			log.Infof("Adding load balancer %+v", s.GetName())

			list = append(list, s)

			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}
