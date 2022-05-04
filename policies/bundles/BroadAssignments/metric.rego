package clouditor.metrics.broad_assignments

import data.clouditor.compare

default applicable = false

default compliant = false

# TODO(lebogg): Not yet in VOC. Check if `rBAC` is correct representation in JSON
broadAssignments := input.rBAC.broadAssignments

applicable {
	broadAssignments
}

compliant {
	# TODO(all): Target value ?
	compare(data.operator, data.target_value, broadAssignments)
}
