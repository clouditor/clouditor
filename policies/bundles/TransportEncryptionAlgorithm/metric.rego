package clouditor.metrics.transport_encryption_algorithm

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

algorithm := enc.algorithm

applicable {
	algorithm != null
}

compliant {
	compare(data.operator, data.target_value, algorithm)
}
