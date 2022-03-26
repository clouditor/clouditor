package clouditor.metrics.identity_password_policy

import data.clouditor.compare
import future.keywords.every
import input as identity

default applicable = false

default compliant = false

applicable {
	# the resource type should be an Identity
	identity.type[_] == "Identity"
}

compliant {
	# we are just assuming, that the standard policy looks good
	compare(data.operator, data.target_value, identity.disablePasswordPolicy)
}
