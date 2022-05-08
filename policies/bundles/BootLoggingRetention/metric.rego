package clouditor.metrics.boot_logging_retention

import data.clouditor.compare

default applicable = false

default compliant = false

retentionPeriod := input.bootLogging.retentionPeriod

applicable {
	retentionPeriod != null
}

compliant {
	compare(data.operator, data.target_value, retentionPeriod)
}
