package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

func (d *azureDiscovery) handleLoadBalancer(lb *armnetwork.LoadBalancer) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: discovery.NewResource(d,
					voc.ResourceID(util.Deref(lb.ID)),
					util.Deref(lb.Name),
					// No creation time available
					nil,
					voc.GeoLocation{
						Region: util.Deref(lb.Location),
					},
					labels(lb.Tags),
					resourceGroupID(lb.ID),
					voc.LoadBalancerType,
					lb,
				),
			},
			Ips:   publicIPAddressFromLoadBalancer(lb),
			Ports: loadBalancerPorts(lb),
		},
		// TODO(all): do we need the httpEndpoint for load balancers?
		HttpEndpoints: []*voc.HttpEndpoint{},
	}
}

// handleApplicationGateway returns the application gateway with its properties
// NOTE: handleApplicationGateway uses the LoadBalancer for now until there is a own resource
func (d *azureDiscovery) handleApplicationGateway(ag *armnetwork.ApplicationGateway) voc.IsNetwork {
	return &voc.LoadBalancer{
		NetworkService: &voc.NetworkService{
			Networking: &voc.Networking{
				Resource: discovery.NewResource(
					d,
					voc.ResourceID(util.Deref(ag.ID)),
					util.Deref(ag.Name),
					nil,
					voc.GeoLocation{Region: util.Deref(ag.Location)},
					labels(ag.Tags),
					resourceGroupID(ag.ID),
					voc.LoadBalancerType,
					ag,
				),
			},
		},
		AccessRestriction: voc.WebApplicationFirewall{
			Enabled: util.Deref(ag.Properties.WebApplicationFirewallConfiguration.Enabled),
		},
	}
}

func (d *azureDiscovery) handleNetworkInterfaces(ni *armnetwork.Interface) voc.IsNetwork {
	return &voc.NetworkInterface{
		Networking: &voc.Networking{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(ni.ID)),
				util.Deref(ni.Name),
				// No creation time available
				nil,
				voc.GeoLocation{
					Region: util.Deref(ni.Location),
				},
				labels(ni.Tags),
				resourceGroupID(ni.ID),
				voc.NetworkInterfaceType,
				ni,
			),
		},
		AccessRestriction: &voc.L3Firewall{
			Enabled: d.nsgFirewallEnabled(ni),
			// Inbound: ,
			// RestrictedPorts: ,
		},
	}
}
