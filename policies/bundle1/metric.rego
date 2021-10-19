package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionEnabled

enc := input.httpEndpoint.transportEncryption

compliant {
	data.operator == "=="
	enc.algorithm == data.target_value
}

applicable {
    input.httpEndpoint[_]
}