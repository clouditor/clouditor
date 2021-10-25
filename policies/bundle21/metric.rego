package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AutomaticUpdatesSecurityOnly

name := "AutomaticUpdatesSecurityOnly"
metricID := 21

autoUpdates := input.automaticUpdates

applicable {
    autoUpdates
}

compliant {
    data.operator == "=="
	autoUpdates.securityOnly == data.target_value
}