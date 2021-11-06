package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingOutput

name := "BootLoggingOutput"
metricData := data.target_value

bootLog := input.bootLog

applicable {
    bootLog
}

compliant {
    isIn(data.target_value, bootLog.output)
}