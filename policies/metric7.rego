package clouditor

default compliant = false

# this is an implementation of metric LoggingEnabled

compliant {
	bootLog := input.bootLog
	bootLog.enabled == true

	# ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'
	OSLog := input.osLog
	OSLog.enabled == true
}