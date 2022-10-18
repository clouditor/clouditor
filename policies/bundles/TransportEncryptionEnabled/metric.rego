package clouditor.metrics.transport_encryption_enabled

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

enabled := input.httpEndpoint.transportEncryption.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
