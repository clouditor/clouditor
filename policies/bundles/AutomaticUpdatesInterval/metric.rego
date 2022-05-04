package clouditor.metrics.automatic_updates_interval

import data.clouditor.compare

default applicable = false

default compliant = false

interval := input.automaticUpdates.interval

applicable {
	interval
}

compliant {
	compare(data.operator, data.target_value, interval)
}
