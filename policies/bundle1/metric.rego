package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionAlgorithm

enc := input.httpEndpoint.transportEncryption
name := "TransportEncryptionAlgorithm"

compliant {
	data.operator == "=="
	enc.algorithm == data.target_value
}

applicable {
    input.httpEndpoint[_]
}