package clouditor

default compliant = false

# this is an implementation of metric TLSVersion

enc := input.httpEndpoint.transportEncryption

compliant {
	data.operator == "=="
	enc.tlsVersion == data.target_value
}

compliant {
	data.operator == ">="
	enc.tlsVersion >= data.target_value
}
