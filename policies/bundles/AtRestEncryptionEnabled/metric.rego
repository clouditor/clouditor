package clouditor.metrics.at_rest_encryption_enabled

import data.clouditor.compare

import input.atRestEncryption as enc

default applicable = false

default compliant = false

applicable {
	enc
}

compliant {
	compare(data.operator, data.target_value, enc[_].enabled)
}
