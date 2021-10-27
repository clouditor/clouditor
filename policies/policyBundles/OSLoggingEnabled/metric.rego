package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingEnabled

name := "OSLoggingEnabled"

OSLog := input.oSLog

applicable {
    OSLog
}

compliant {
    data.operator == "=="
	OSLog.enabled == data.target_value
}