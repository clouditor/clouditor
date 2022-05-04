package clouditor.metrics.l_3_firewall_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.l3Firewall.enabled

applicable {
	enabled
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
