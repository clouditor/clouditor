package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingRetention

name := "BootLoggingRetention"

bootLog := input.bootLog

applicable {
    bootLog
}

compliant {
    compare(data.operator, data.target_value, bootLog.retentionPeriod)
}