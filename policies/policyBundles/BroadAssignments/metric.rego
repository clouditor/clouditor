package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BroadAssignments

name := "BroadAssignments"

# TODO(lebogg): Check if `rBAC` is correct representation in JSON
rbac := input.rBAC

applicable {
    rbac
}

compliant {
    data.operator == "<="
    # TODO(all): Target value ?
	rbac.broadAssignments <= data.target_value
}