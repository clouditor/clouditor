package clouditor.metrics.l_3_firewall_restricted_ports

import data.clouditor.compare

default applicable = false

default compliant = false

restrictedPorts := input.l3Firewall.restrictedPorts

applicable {
	restrictedPorts
}

# TODO(all): Maybe change restrictet ports to array of strings. See comment in Ontology.
compliant {
	compare(data.operator, data.target_value, restrictedPorts)
}
