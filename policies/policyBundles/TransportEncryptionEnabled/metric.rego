package clouditor

default compliant = false
default applicable = false

# this is an implementation of metric TransportEncryptionEnabled

name := "TransportEncryptionEnabled"

endpoint := input.httpEndpoint

applicable {
    endpoint
}

# TODO(all): Alternatively, curly braces can be removed and a single assignment used. But for readability and consistency (having multiple compares, see mutual auth) I let it this way?
compliant {
    compare(data.operator, data.target_value, endpoint.transportEncryption.enabled)
}