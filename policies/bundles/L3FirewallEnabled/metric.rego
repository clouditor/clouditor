package clouditor.metrics.l_3_firewall_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.accessRestriction.enabled

applicable {
	enabled != null
}

compliant {
	compare("isIn",  "NetworkInterface", input.type)
	compare(data.operator, data.target_value, enabled)
}
