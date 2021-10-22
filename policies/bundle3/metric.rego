package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionEnabled

name := "TransportEncryptionEnabled"
metricID := 3

endpoint := input.httpEndpoint

applicable {
    endpoint
}

compliant {
	data.operator == "=="
	endpoint.transportEncryption.enabled == data.target_value
}