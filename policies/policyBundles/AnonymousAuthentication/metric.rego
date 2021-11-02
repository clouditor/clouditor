package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnonymousAuthentication

name := "AnonymousAuthentication"

cba := input.certificateBasedAuthentication

applicable {
    cba
}

compliant {
    data.operator == "=="
	cba.enabled == data.target_value
}