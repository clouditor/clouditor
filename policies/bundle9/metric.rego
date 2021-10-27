package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingOutput

name := "BootLoggingOutput"
metricID := 9

bootLog := input.bootLog

applicable {
    bootLog
}

compliant {
    data.operator == "=="
    some i
    some j
	bootLog.output[i] == data.target_value[j]
	bootLog.output[j] == data.target_value[i]
}