package clouditor.metrics.single_sign_on_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.singleSignOn.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
