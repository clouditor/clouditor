package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingOutput

name := "OSLoggingOutput"
metricData := data.target_value

OSLog := input.oSLog

applicable {
    OSLog
}

compliant {
    data.operator == "=="
    # Current implementation: It is enough that one output is one of target_values
	OSLog.output[_] == data.target_value[_]
}