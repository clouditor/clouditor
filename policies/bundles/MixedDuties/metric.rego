package clouditor.metrics.mixed_duties

import data.clouditor.compare

default applicable = false

default compliant = false

# TODO(lebogg): Check if `rBAC` is correct representation in JSON
mixedDuties := input.rBAC.mixedDuties

applicable {
	mixedDuties
}

compliant {
	# TODO(all): Target value ?
	compare(data.operator, data.target_value, mixedDuties)
}
