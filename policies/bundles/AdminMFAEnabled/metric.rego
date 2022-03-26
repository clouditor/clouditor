package clouditor.metrics.admin_mfa_enabled

import future.keywords.every
import data.clouditor.compare
import input as account

default applicable = false

default compliant = false

applicable {
    # we are only interested in some kind of admin user    
	account.isAdmin
}

compliant {
    # count the number of "factors"
	compare(data.operator, data.target_value, account.authenticity)

    # also make sure, that we do not have any "NoAuthentication" in the factor and all are activated
    every factor in account.authenticity {
        # TODO(oxisto): we do not have this type property (yet)
        not factor.type == "NoAuthentication"

        factor.activated == true
    }
}
