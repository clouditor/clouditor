package clouditor.metrics.os_logging_retention

import data.clouditor.compare

default applicable = false

default compliant = false

retentionPeriod := input.oSLogging.retentionPeriod

applicable {
	retentionPeriod != null
}

compliant {
	compare(data.operator, data.target_value, retentionPeriod)
}
