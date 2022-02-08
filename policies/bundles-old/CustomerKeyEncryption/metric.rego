package clouditor.customer.key.encryption

default applicable = false

default compliant = false

enc := input.atRestEncryption

applicable {
	enc
}

compliant {
	# Check if keyUrl is set. It is only set in the customer key encryption case.
	enc.keyUrl
}
