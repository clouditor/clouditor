package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AutomaticUpdatesEnabled

name := "AutomaticUpdatesEnabled"

autoUpdates := input.automaticUpdates

applicable {
    autoUpdates
}

compliant {
    compare(data.operator, data.target_value, autoUpdates.enabled)
}