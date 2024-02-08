package clouditor.metrics.os_logging_output

import data.clouditor.compare

# spelling is incorrect, need to change after we fix it in owl2proto
import input.oslogging as logging

#import input.osLogging as logging

default applicable = false

default compliant = false

metricConfiguration := data.target_value

applicable {
	logging
}

compliant {
	compare(data.operator, data.target_value, count(logging.loggingServiceIds))
}
