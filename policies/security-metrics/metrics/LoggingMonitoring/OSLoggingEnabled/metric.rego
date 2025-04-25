package metrics.logging_monitoring.os_logging_enabled

import data.compare
import rego.v1
import input.osLogging as logging

default applicable = false

default compliant = false

applicable if {
	logging
}

compliant if {
	compare(data.operator, data.target_value, logging.enabled)
}
