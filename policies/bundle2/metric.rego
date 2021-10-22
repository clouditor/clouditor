package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TLSVersion

name := "TLSVersion"
metricID := 2

endpoint := input.httpEndpoint

applicable {
    endpoint
}

compliant {
	data.operator == "=="
	endpoint.transportEncryption.tlsVersion == data.target_value
}

compliant {
	data.operator == ">="
	endpoint.transportEncryption.tlsVersion >= data.target_value
}