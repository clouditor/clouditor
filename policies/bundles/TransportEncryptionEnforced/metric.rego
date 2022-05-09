package clouditor.metrics.transport_encryption_enforced

import data.clouditor.compare

default compliant = false

default applicable = false

enforced := input.httpEndpoint.transportEncryption.enforced

applicable {
	enforced != null
}

compliant {
	compare(data.operator, data.target_value, enforced)
}
