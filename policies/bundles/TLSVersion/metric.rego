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
    # If target_value is a list of strings/numbers
	isIn(data.target_value, endpoint.transportEncryption.tlsVersion)
}{
    # If target_value is the version number represented as int/float
    compare(data.operator, data.target_value, endpoint.transportEncryption.tlsVersion)
}