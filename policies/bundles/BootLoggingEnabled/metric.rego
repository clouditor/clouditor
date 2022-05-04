package clouditor.metrics.boot_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.bootLogging.enabled

applicable {
	enabled
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
