package clouditor.metrics.automatic_updates_security_only

import data.clouditor.compare

default applicable = false

default compliant = false

autoUpdates := input.automaticUpdates

applicable {
	autoUpdates
}

compliant {
	compare(data.operator, data.target_value, autoUpdates.securityOnly)
}
