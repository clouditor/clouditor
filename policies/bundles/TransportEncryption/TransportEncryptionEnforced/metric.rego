package clouditor.metrics.transport_encryption_enforced

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

applicable if {
	enc
}

compliant if {
	compare(data.operator, data.target_value, enc.enforced)
}
