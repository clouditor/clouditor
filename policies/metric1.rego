package clouditor

default compliant = false

# this is an implementation of metric 1 (Transport Encryption)

compliant {
	tls := input.httpEndpoint.transportEncryption
	tls.enabled == true
}
