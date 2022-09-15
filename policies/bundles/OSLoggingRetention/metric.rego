package clouditor.metrics.os_logging_retention

import data.clouditor.compare

default applicable = false

default compliant = false

retentionPeriod := input.osLogging.retentionPeriod

applicable {
	retentionPeriod != null
}

compliant {
	# time.Duration is nanoseconds, we want to convert this to hours
	days := retentionPeriod / (((1000 * 1000) * 1000) * 3600)

	compare(data.operator, data.target_value, days)
}
