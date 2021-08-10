package clouditor

default compliant = false

compliant {
	tls := input.httpEndpoint.transportEncryption
	tls.algorithm == true
	goodVersion(tls)
}

goodVersion(tls) {
	tls.version == "TLS 1.2"
}

goodVersion(tls) {
	tls.version == "TLS 1.3"
}
