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
    compare(data.operator, data.target_value, OSLog.enabled)
}