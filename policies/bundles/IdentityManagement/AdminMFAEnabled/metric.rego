package clouditor.metrics.admin_mfa_enabled

import data.clouditor.compare
import future.keywords.every
import input as identity

default applicable = false

default compliant = false

applicable {
	# we are only interested in some kind of privileged user    
	identity.privileged
}

compliant {
	# count the number of "factors"
	compare(data.operator, data.target_value, count(identity.authenticity))

	# also make sure, that we do not have any "NoAuthentication" in the factor and all are activated
	every factor in identity.authenticity {
		# TODO(oxisto): we do not have this type property (yet)
		not factor.type == "NoAuthentication"

		factor.activated == true
	}
}
