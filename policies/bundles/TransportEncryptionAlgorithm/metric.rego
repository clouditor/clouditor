package clouditor.metrics.transport_encryption_algorithm

import data.clouditor.compare
import input.httpEndpoint as endpoint

default compliant = false

default applicable = false

algorithm := endpoint.transportEncryption.algorithm

applicable {
	algorithm != null
}

compliant {
	compare(data.operator, data.target_value, algorithm)
}
