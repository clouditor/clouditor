package metrics.automatic_updates_enabled

import data.compare
import input.automaticUpdates as au

default applicable = false

default compliant = false

applicable if {
	au
}

compliant if {
	compare(data.operator, data.target_value, au.enabled)
}
