package clouditor.metrics.os_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

OSLogging := input.oSLogging

applicable {
	OSLogging
}

compliant {
	compare(data.operator, data.target_value, OSLogging.enabled)
}
