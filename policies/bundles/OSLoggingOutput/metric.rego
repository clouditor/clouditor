package clouditor

default applicable = false

default compliant = false

metricData := data.target_value

OSLog := input.oSLog

applicable {
	OSLog
}

compliant {
	isIn(data.target_value, OSLog.output)
}
