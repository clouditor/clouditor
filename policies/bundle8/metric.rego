package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingRetention

name := "BootLoggingRetention"
metricID := 8

bootLog := input.bootLog

applicable {
    bootLog
}

compliant {
    data.operator == ">="
	bootLog.retentionPeriod >= data.target_value
}