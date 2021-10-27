package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionEnforced

name := "TransportEncryptionEnforced"

endpoint := input.httpEndpoint

applicable {
    endpoint
}

compliant {
	data.operator == "=="
	endpoint.transportEncryption.enforced == data.target_value
}