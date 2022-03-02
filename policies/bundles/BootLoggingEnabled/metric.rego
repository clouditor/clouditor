package clouditor.metrics.boot_logging_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

bootLog := input.bootLog

applicable {
	bootLog
}

compliant {
	compare(data.operator, data.target_value, bootLog.enabled)
}
