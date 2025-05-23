package clouditor.metrics.automatic_updates_security_only

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

applicable if {
	am
}

compliant if {
	compare(data.operator, data.target_value, am.securityOnly)
}
