package clouditor.automatic.updates.interval

import data.clouditor.compare

default applicable = false

default compliant = false

autoUpdates := input.automaticUpdates

applicable {
	autoUpdates
}

compliant {
	compare(data.operator, data.target_value, autoUpdates.interval)
}
