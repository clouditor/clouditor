package ionos

import (
	"clouditor.io/clouditor/v2/internal/util"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func (d *ionosDiscovery) getRestrictedPorts(nic ionoscloud.Nic) []string {
	var (
		restrictedPortsList []string
	)

	if nic.Entities == nil || nic.Entities.Firewallrules == nil || nic.Entities.Firewallrules.Items == nil {
		return restrictedPortsList
	}

	for _, rule := range *nic.Entities.Firewallrules.Items {
		if rule.Properties == nil {
			continue
		}

		if util.Deref(rule.Properties.PortRangeStart) == 0 && util.Deref(rule.Properties.PortRangeEnd) == 0 {
			// If no port range is specified, it means all ports are allowed
			restrictedPortsList = append(restrictedPortsList, "all")
		} else if rule.Properties.PortRangeStart == rule.Properties.PortRangeEnd {
			// If the port range is a single port, add that port to the list
			restrictedPortsList = append(restrictedPortsList, string(util.Deref(rule.Properties.PortRangeStart)))
		} else if rule.Properties.PortRangeStart != nil && rule.Properties.PortRangeEnd != nil {
			// If the port range is specified, add each port in the range to the list
			for port := util.Deref(rule.Properties.PortRangeStart); port <= util.Deref(rule.Properties.PortRangeEnd); port++ {
				restrictedPortsList = append(restrictedPortsList, string(port))
			}
		}

	}

	return restrictedPortsList
}
