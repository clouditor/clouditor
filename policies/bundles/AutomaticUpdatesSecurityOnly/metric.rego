package clouditor.metrics.automatic_updates_security_only

import data.clouditor.compare

default applicable = false

default compliant = false

securityOnly := input.automaticUpdates.securityOnly

applicable {
	securityOnly
}

compliant {
	compare(data.operator, data.target_value, securityOnly)
}
