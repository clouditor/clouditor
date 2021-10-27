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
	data.operator == "=="
	endpoint.transportEncryption.algorithm == data.target_value
}

