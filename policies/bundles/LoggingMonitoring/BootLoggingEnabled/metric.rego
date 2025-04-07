package clouditor.metrics.boot_logging_enabled

import data.clouditor.compare
import input.bootLogging as logging

default applicable = false

default compliant = false

applicable if {
	logging
}

compliant if {
	compare(data.operator, data.target_value, logging.enabled)
}
