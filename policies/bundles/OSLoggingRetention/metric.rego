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
    compare(data.operator, data.target_value, OSLog.retentionPeriod)
}