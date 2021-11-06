package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingEnabled

name := "BootLoggingEnabled"

bootLog := input.bootLog

applicable {
    bootLog
}

compliant {
    compare(data.operator, data.target_value, bootLog.enabled)
}