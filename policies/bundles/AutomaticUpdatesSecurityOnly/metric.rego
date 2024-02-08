package clouditor.metrics.automatic_updates_security_only

import data.clouditor.compare
import input.automaticSecurityUpdates as am

default applicable = false

default compliant = false

applicable {
	am
}

compliant {
	compare(data.operator, data.target_value, am.securityOnly)
}
