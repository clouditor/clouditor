package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AutomaticUpdatesSecurityOnly

name := "AutomaticUpdatesSecurityOnly"

autoUpdates := input.automaticUpdates

applicable {
    autoUpdates
}

compliant {
    compare(data.operator, data.target_value, autoUpdates.securityOnly)
}