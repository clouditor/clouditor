package clouditor.metrics.l_3_firewall_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.accessRestriction.enabled

applicable if {
	enabled != null
	compare("isIn",  "NetworkInterface", input.type)
}

compliant if {
	compare(data.operator, data.target_value, enabled)
}
