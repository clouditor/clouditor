package clouditor.metrics.os_logging_retention

import data.clouditor.compare

import input.osLogging as logging

default applicable = false

default compliant = false

applicable {
	logging
}

compliant {
	# time.Duration is nanoseconds, we want to convert this to hours
	days := time.parse_duration_ns(logging.retentionPeriod) / (((1000 * 1000) * 1000) * 3600)

	compare(data.operator, data.target_value, days)
}
