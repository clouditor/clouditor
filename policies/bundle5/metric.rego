package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric EncryptionAtRestEnabled

name := "EncryptionAtRestEnabled"

applicable {
    input.atRestEncryption[_]
}

compliant {
	enc := input.atRestEncryption
    data.operator == "=="
	enc.enabled == data.target_value
}
