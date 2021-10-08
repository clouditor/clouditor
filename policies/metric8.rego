package clouditor

default compliant = false

# this is an implementation of metric LogRetention

compliant {
	bootLog := input.bootLog
	bootLog.retention >= 90

	# ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'
	OSLog := input.osLog
	OSLog.retention >= 90
}