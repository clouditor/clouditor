package clouditor.metrics.os_logging_output

import data.clouditor.compare

default applicable = false

default compliant = false

metricConfiguration := data.target_value

output := input.oSLogging.loggingService

applicable {
	output != null
}

compliant {
	compare(data.operator, data.target_value, output)
}
