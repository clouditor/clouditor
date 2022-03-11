package clouditor.metrics.boot_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

bootLogging := input.bootLogging

applicable {
	bootLogging
}

compliant {
	compare(data.operator, data.target_value, bootLogging.enabled)
}
