package clouditor.metrics.admin_mixed_duties

import data.clouditor.compare
import future.keywords.every
import input as account

default applicable = false

default compliant = false

applicable {
	# we are only interested in some kind of privileged user    
	account.privileged
}

compliant {
	compare(data.operator, data.target_value, account.authorization.mixedDuties)
}
