package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingEnabled

OSLog := input.OSLog

applicable {
    OSLog[_]
}

compliant {
	# ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'
    data.operator == "=="
	OSLog.enabled == data.target_value
}