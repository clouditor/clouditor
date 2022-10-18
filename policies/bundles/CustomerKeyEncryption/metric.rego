package clouditor.metrics.customer_key_encryption

default applicable = false

default compliant = false

enc := input.atRestEncryption

applicable {
	enc != null
}

compliant {
	# Check if keyUrl is set (not an empty string)
	enc.keyUrl != ""
}
