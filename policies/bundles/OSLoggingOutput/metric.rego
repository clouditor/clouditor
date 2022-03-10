package clouditor.metrics.os_logging_output

import data.clouditor.isIn

default applicable = false

default compliant = false

metricConfiguration := data.target_value

OSLogging := input.oSLogging.logging.loggingService

applicable {
	OSLogging
}

compliant {
	isIn(data.target_value, OSLogging.output)
}
