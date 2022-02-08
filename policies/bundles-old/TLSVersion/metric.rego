package clouditor.tls.version

import data.clouditor.compare
import data.clouditor.isIn

default compliant = false

default applicable = false

endpoint := input.httpEndpoint

applicable {
	endpoint
}

compliant {
	# If target_value is a list of strings/numbers
	isIn(data.target_value, endpoint.transportEncryption.tlsVersion)
}

compliant {
	# If target_value is the version number represented as int/float
	compare(data.operator, data.target_value, endpoint.transportEncryption.tlsVersion)
}
