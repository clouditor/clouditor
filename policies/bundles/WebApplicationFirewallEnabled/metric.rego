package clouditor.metrics.web_application_firewall_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.accessRestriction.enabled

applicable {
	enabled != null
	compare("isIn",  "LoadBalancer", input.type)
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
