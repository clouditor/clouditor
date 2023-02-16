package clouditor.metrics.tls_version

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

applicable {
	enc
}

compliant {
	# If target_value is the version number represented as int/float
	compare(data.operator, data.target_value, enc.tlsVersion)
}
