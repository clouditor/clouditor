package clouditor.metrics.anonymous_authentication

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.certificateBasedAuthentication.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
