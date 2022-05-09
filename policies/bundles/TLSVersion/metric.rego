package clouditor.metrics.tls_version

import data.clouditor.compare
import data.clouditor.isIn

default compliant = false

default applicable = false

version := input.httpEndpoint.transportEncryption.tlsVersion

applicable {
	version != null
}

compliant {
	# If target_value is a list of strings/numbers
	isIn(data.target_value, version)
}

compliant {
	# If target_value is the version number represented as int/float
	compare(data.operator, data.target_value, version)
}
