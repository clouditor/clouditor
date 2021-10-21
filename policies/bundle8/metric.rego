package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingRetention

name := "BootLoggingRetention"

bootLog := input.bootLog

applicable {
    bootLog[_]
}

compliant {
    data.operator == ">="
	bootLog.retentionPeriod >= data.target_value
}