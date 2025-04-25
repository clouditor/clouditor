package metrics.logging_monitoring.os_logging_retention

import data.compare
import rego.v1
import input.osLogging as logging

default applicable = false

default compliant = false

applicable if {
	logging
}

compliant if {
	# time.Duration is nanoseconds, we want to convert this to days
	days := time.parse_duration_ns(logging.retentionPeriod) / (((1000 * 1000) * 1000) * 3600)

	compare(data.operator, data.target_value, days)
}
