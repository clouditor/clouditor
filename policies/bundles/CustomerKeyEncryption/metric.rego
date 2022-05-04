package clouditor.metrics.customer_key_encryption

default applicable = false

default compliant = false

keyUrl := input.atRestEncryption.keyUrl

applicable {
	keyUrl
}

compliant {
	# Check if keyUrl is set (not an empty string)
	keyUrl != ""
}
