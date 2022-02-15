package clouditor.metrics.transport_encryption_algorithm

import data.clouditor.compare
import input.httpEndpoint as endpoint

default compliant = false

default applicable = false

applicable {
	endpoint
}

compliant {
	compare(data.operator, data.target_value, endpoint.transportEncryption.algorithm)
}
