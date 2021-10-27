package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AutomaticUpdatesInterval

name := "AutomaticUpdatesInterval"
metricID := 20

autoUpdates := input.automaticUpdates

applicable {
    autoUpdates
}

compliant {
    data.operator == "=="
	autoUpdates.interval == data.target_value
}