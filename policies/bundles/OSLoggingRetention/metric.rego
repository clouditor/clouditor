package clouditor.metrics.os_logging_retention

import data.clouditor.compare

default applicable = false

default compliant = false

OSLogging := input.osLogging

applicable {
	OSLogging
}

compliant {
	# time.Duration is nanoseconds, we want to convert this to hours
	days := OSLogging.retentionPeriod / (1000*1000*1000*3600)
	
	compare(data.operator, data.target_value, days)
}
