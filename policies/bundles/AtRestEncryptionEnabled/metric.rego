package clouditor

default applicable = false

default compliant = false

enc := input.atRestEncryption

applicable {
	enc
}

compliant {
	compare(data.operator, data.target_value, enc.enabled)
}
