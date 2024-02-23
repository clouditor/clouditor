package clouditor.metrics.transport_encryption_enforced

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

applicable {
	enc != null
}

compliant {
	compare(data.operator, data.target_value, enc.enforced)
}
