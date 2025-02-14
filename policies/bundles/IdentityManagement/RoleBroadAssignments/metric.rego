package clouditor.metrics.role_broad_assignments

import data.clouditor.compare

# TODO(lebogg): Not yet in VOC. Check if `rBAC` is correct representation in JSON
import input.rBAC as rbac

default applicable = false

default compliant = false

applicable if {
	rbac.broadAssignments
}

compliant if {
	# TODO(all): Target value ?
	compare(data.operator, data.target_value, rbac.broadAssignments)
}
