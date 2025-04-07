package clouditor.metrics.anonymous_authentication

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.certificateBasedAuthentication.enabled

applicable if {
	enabled != null
}

compliant if {
	compare(data.operator, data.target_value, enabled)
}
