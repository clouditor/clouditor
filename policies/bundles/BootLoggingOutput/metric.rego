package clouditor.metrics.boot_logging_output

import data.clouditor.isIn

default applicable = false

default compliant = false

metricConfiguration := data.target_value

bootLogging := input.bootLogging

applicable {
	bootLogging
}

compliant {
	isIn(data.target_value, bootLogging.loggingService)
}
