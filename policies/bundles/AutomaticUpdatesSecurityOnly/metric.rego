package clouditor

default applicable = false

default compliant = false

autoUpdates := input.automaticUpdates

applicable {
	autoUpdates
}

compliant {
	compare(data.operator, data.target_value, autoUpdates.securityOnly)
}
