package clouditor.metrics.os_logging_output

import data.clouditor.compare
import input.oSLogging as logging

default applicable = false

default compliant = false

metricConfiguration := data.target_value

applicable {
	logging
}

compliant {
	compare(data.operator, data.target_value, count(logging.loggingService))
}
