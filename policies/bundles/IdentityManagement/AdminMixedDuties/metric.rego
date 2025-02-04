package clouditor.metrics.admin_mixed_duties

import data.clouditor.compare
import future.keywords.every
import input as identity

default applicable = false

default compliant = false

applicable if {
	# we are only interested in some kind of privileged user    
	identity.privileged
}

compliant if {
	compare(data.operator, data.target_value, identity.authorization.mixedDuties)
}
