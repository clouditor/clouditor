package clouditor

default compliant = false

# this is an implementation of metric TransportEncryptionEnforced

enc := input.httpEndpoint.transportEncryption

compliant {
	data.operator == "=="
	enc.enforced == data.target_value
}
