package clouditor.metrics.anonymous_authentication

import data.clouditor.compare

default applicable = false

default compliant = false

cba := input.certificateBasedAuthentication

applicable {
	cba
}

compliant {
	compare(data.operator, data.target_value, cba.enabled)
}
