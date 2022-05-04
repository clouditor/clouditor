package clouditor.metrics.web_application_firewall_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.webApplicationFirewall.enabled

applicable {
	enabled
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
