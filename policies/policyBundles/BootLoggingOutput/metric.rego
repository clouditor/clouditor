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
    data.operator == "=="
    # Current implementation: It is enough that one output is one of target_values
	bootLog.output[_] == data.target_value[_]
}