package clouditor.metrics.boot_logging_output

import data.clouditor.isIn

default applicable = false

default compliant = false

metricConfiguration := data.target_value

bootLog := input.bootLog

applicable {
	bootLog
}

compliant {
	isIn(data.target_value, bootLog.output)
}
