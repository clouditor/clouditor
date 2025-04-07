package clouditor.metrics.customer_key_encryption

default applicable = false

default compliant = false

import input.atRestEncryption.customerKeyEncryption as cke

applicable if {
	cke
}

compliant if {
	# Check if keyUrl is set (not an empty string)
	cke.keyUrl != ""
}
