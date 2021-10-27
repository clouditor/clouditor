package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TLSVersion

name := "TLSVersion"

endpoint := input.httpEndpoint

applicable {
    endpoint
}

compliant {
	data.operator == "=="
	# Assess version with each compliant version in data
	endpoint.transportEncryption.tlsVersion == data.target_value[_]
}

# Currently not working since version is string in ontology
compliant {
	data.operator == ">="
	endpoint.transportEncryption.tlsVersion >= data.target_value
}