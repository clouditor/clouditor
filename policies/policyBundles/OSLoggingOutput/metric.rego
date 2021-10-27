package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingOutput

name := "OSLoggingOutput"

OSLog := input.OSLog

applicable {
    OSLog
}

compliant {
    data.operator == "=="
	OSLog.output == data.target_value
}