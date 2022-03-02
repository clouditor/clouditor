package clouditor.metrics.transport_encryption_enforced

import data.clouditor.compare

default compliant = false

default applicable = false

endpoint := input.httpEndpoint

applicable {
	endpoint
}

compliant {
	compare(data.operator, data.target_value, endpoint.transportEncryption.enforced)
}
