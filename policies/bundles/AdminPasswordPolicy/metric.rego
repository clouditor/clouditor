package clouditor.metrics.admin_password_policy

import future.keywords.every
import data.clouditor.compare
import input as account

default applicable = false

default compliant = false

applicable {
    # the resource type should be an account
	account.type[_] == "Account"
}

compliant {
    # we are just assuming, that the standard policy looks good
    compare(data.operator, data.target_value, account.disablePasswordPolicy)
}
