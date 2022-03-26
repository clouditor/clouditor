package clouditor.metrics.account_recent_activity

import data.clouditor.compare
import future.keywords.every
import input as account

default applicable = false

default compliant = false

applicable {
	# we are only interested in some kind of admin user    
	account.isAdmin

	# and we are also only interested in active accounts
	account.activated
}

compliant {
	ts := time.parse_rfc3339_ns(account.lastActivity)
	now := time.now_ns()

	#window := ((((90 * 24) * 3600) * 1000) * 1000) * 1000
	window := ((((data.target_value * 24) * 3600) * 1000) * 1000) * 1000

	#now - ts <= window
	compare(data.operator, window, now - ts)
}
