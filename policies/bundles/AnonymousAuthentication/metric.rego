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
    compare(data.operator, data.target_value, cba.enabled)
}