package clouditor.metrics.os_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.oSLogging.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
