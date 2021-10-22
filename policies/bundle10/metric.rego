package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingEnabled

name := "OSLoggingEnabled"
metricID := 10

OSLog := input.OSLog

applicable {
    OSLog
}

compliant {
	# ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'
    data.operator == "=="
	OSLog.enabled == data.target_value
}