package clouditor.metrics.automatic_updates_interval

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

applicable {
	am
}

compliant {
	# time.Duration is hours, we want to convert this to days
	compare(data.operator, data.target_value, am.interval / (((1000 * 1000) * 1000) * 3600))
}
