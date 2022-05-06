package clouditor.metrics.automatic_updates_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

enabled := input.automaticUpdates.enabled

applicable {
	enabled != null
}

compliant {
	compare(data.operator, data.target_value, enabled)
}
