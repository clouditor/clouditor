package clouditor.metrics.boot_logging_enabled

import data.clouditor.compare
import input.bootLogging as logging

default applicable = false

default compliant = false

applicable {
	logging
}

compliant {
	compare(data.operator, data.target_value, logging.enabled)
}
