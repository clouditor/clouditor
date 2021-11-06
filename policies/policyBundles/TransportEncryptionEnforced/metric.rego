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
    compare(data.operator, data.target_value, endpoint.transportEncryption.enforced)
}