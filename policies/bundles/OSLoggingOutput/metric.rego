package clouditor.metrics.os_logging_output

import data.clouditor.isIn

default applicable = false

default compliant = false

metricConfiguration := data.target_value

output := input.oSLogging.loggingService

applicable {
	output != null
}

compliant {
	isIn(data.target_value, output)
}
