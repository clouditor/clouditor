package clouditor

default applicable = false

default compliant = false

# TODO(lebogg): Not yet in VOC. Check if `rBAC` is correct representation in JSON
rbac := input.rBAC

applicable {
	rbac
}

compliant {
	# TODO(all): Target value ?
	compare(data.operator, data.target_value, rbac.broadAssignments)
}
