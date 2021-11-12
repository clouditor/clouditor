package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric MixedDuties

name := "MixedDuties"

# TODO(lebogg): Check if `rBAC` is correct representation in JSON
rbac := input.rBAC

applicable {
    rbac
}

compliant {
    # TODO(all): Target value ?
    compare(data.operator, data.target_value, rbac.mixedDuties)
}