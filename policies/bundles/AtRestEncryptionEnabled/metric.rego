package clouditor.metrics.at_rest_encryption_enabled

import data.clouditor.compare
import input.atRestEncryption as enc

default applicable = false

default compliant = false

enabled := enc.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
