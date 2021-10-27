package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingRetention

name := "OSLoggingRetention"

OSLog := input.oSLog

applicable {
    OSLog
}

compliant {
    data.operator == ">="
	OSLog.retentionPeriod >= 35
}