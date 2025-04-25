package metrics.iam.identity_password_policy_enabled

import data.compare
import rego.v1
import input as identity

default applicable = false

default compliant = false

applicable if {
	# the resource type should be an Identity
	identity.type[_] == "Identity"
}

compliant if {
	# we are just assuming that the standard policy looks good
	compare(data.operator, data.target_value, identity.disablePasswordPolicy)
}
