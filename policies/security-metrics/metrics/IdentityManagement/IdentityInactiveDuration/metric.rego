package metrics.iam.identity_inactive_duration

import data.compare
import rego.v1
import input as identity

default applicable = false

default compliant = false

applicable if {
    # we are only interested in active accounts, deactivated accounts are already inactive
	identity.activated
}

compliant if {
	ts := time.parse_rfc3339_ns(identity.lastActivity)
	now := time.now_ns()

	#window := ((((90 * 24) * 3600) * 1000) * 1000) * 1000
	window := ((((data.target_value * 24) * 3600) * 1000) * 1000) * 1000

	#now - ts <= window
	compare(data.operator, now - ts, window)
}
