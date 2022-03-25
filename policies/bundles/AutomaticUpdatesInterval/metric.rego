package clouditor.metrics.automatic_updates_interval

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

applicable {
	am
}

tv := data.target_value

compliant {
	# time.Duration is nanoseconds, we want to convert this to hours
	compare(data.operator, data.target_value, am.interval / (1000*1000*1000*3600))
}
