package clouditor.metrics.customer_key_encryption

default applicable = false

default compliant = false

import input.atRestEncryption as enc

applicable {
	enc
}

compliant {
	# Check if keyUrl is set (not an empty string)
	enc.keyUrl != ""
}
