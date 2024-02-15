package clouditor.metrics.os_logging_enabled

import data.clouditor.compare

import input.osLogging as logging

default applicable = false

default compliant = false

applicable {
	logging
}

compliant {
	compare(data.operator, data.target_value, logging.enabled)
}
