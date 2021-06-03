package clouditor

default compliant = false

compliant {
	tls := input.httpEndpoint.transportEncryption
	tls.enabled == true
	tls.tlsVersion == "TLS1_2"
}
