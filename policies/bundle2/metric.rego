package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TLSVersion

enc := input.httpEndpoint.transportEncryption
name := "TLSVersion"

compliant {
	data.operator == "=="
	enc.tlsVersion == data.target_value
}

compliant {
	data.operator == ">="
	enc.tlsVersion >= data.target_value
}

applicable {
    input.httpEndpoint[_]
}