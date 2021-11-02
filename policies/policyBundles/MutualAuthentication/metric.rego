package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric MutualAuthentication

name := "MutualAuthentication"

cba := input.certificateBasedAuthentication
enc := input.httpEndpoint.transportEncryption

applicable {
    cba
}

compliant {
    data.operator == "=="
	cba.enabled == data.target_value
	enc.enabled == data.target_value
}