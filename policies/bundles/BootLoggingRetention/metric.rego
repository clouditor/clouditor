package clouditor

default applicable = false

default compliant = false

bootLog := input.bootLog

applicable {
	bootLog
}

compliant {
	compare(data.operator, data.target_value, bootLog.retentionPeriod)
}
