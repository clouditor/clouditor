package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionAlgorithm

name := "TransportEncryptionAlgorithm"

endpoint := input.httpEndpoint

applicable {
    endpoint
}

compliant {
    compare(data.operator, data.target_value, endpoint.transportEncryption.algorithm)
}

