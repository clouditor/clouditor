package clouditor

default compliant = false

# this is an implementation of metric TransportEncryptionEnabled

enc := input.httpEndpoint.transportEncryption

compliant {
	data.operator == "=="
	enc.algorithm == data.target_value
}
