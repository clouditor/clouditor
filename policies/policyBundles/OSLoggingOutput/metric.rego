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
    isIn(data.target_value, OSLog.output)
}