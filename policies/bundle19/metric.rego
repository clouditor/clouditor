package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AutomaticUpdatesEnabled

name := "AutomaticUpdatesEnabled"
metricID := 19

autoUpdates := input.automaticUpdates

applicable {
    autoUpdates
}

compliant {
    data.operator == "<="
	autoUpdates.enabled == data.target_value
}