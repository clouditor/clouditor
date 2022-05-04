package clouditor.metrics.transport_encryption_enabled

import data.clouditor.compare

default compliant = false

default applicable = false

enabled := input.httpEndpoint.transportEncryption.enabled

applicable {
	enabled
}

# TODO(all): Alternatively, curly braces can be removed and a single assignment used. But for readability and consistency (having multiple compares, see mutual auth) I let it this way?
compliant {
	compare(data.operator, data.target_value, enabled)
}
