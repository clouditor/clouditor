package clouditor.metrics.at_rest_encryption_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.atRestEncryption.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
