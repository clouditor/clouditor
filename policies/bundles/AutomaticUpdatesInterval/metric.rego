package clouditor.metrics.automatic_updates_interval

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

applicable {
	am
}

compliant {
	# time.Duration is nanoseconds, we want to convert this to days
	compare(data.operator, data.target_value, time.parse_duration_ns(am.interval) / (1000000000 * 86400)) # nanoseconds to seconds (/1000000000), seconds to days (/(60*60*24) = 86400) 
}
