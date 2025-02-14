package clouditor.metrics.automatic_updates_interval

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

applicable if {
	am
}

compliant if {
	# Check if interval is > 0.
	# The discoverer sets the interval to 0 if the the automatic update is not enabled. If we do not check 'interval > 0' it can result in 'AutomaticUpdatesEnabled=false' and  'AutomaticUpdatesInterval=true'. 
	compare(">", 0, time.parse_duration_ns(am.interval) / (1000000000 * 86400))
	# time.Duration is nanoseconds, we want to convert this to days
	compare(data.operator, data.target_value, time.parse_duration_ns(am.interval) / (1000000000 * 86400)) # nanoseconds to seconds (/1000000000), seconds to days (/(60*60*24) = 86400) 
}
