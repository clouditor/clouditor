package clouditor

default compliant = false

# this is an implementation of metric TransportEncryptionEnforced

compliant {
	tls := input.httpEndpoint.transportEncryption
	tls.enforced == true
}