package clouditor.transport.encryption.algorithm

import data.clouditor.compare

default compliant = false

default applicable = false

endpoint := input.httpEndpoint

applicable {
	endpoint
}

compliant {
	compare(data.operator, data.target_value, endpoint.transportEncryption.algorithm)
}
