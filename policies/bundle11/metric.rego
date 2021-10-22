package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingRetention

name := "OSLoggingRetention"
metricID := 11

OSLog := input.OSLog

applicable {
    OSLog
}

compliant {
	# ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'
    data.operator == ">="
	OSLog.retentionPeriod >= 35
}