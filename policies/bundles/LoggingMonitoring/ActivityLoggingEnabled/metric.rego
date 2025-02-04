package clouditor.metrics.activity_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.resourceLogging.enabled

applicable if {
	enabled != null
}

compliant if {
	compare(data.operator, data.target_value, enabled)
}