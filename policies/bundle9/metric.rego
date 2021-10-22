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
    data.operator == ">="
	bootLog.output == data.target_value
}