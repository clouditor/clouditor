package clouditor.metrics.boot_logging_output

import data.clouditor.isIn
import input.bootLogging as logging

default applicable = false

default compliant = false

metricConfiguration := data.target_value

applicable {
	logging
}

compliant {
	isIn(data.target_value, logging.loggingService)
}
