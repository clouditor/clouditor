package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric EncryptionAtRestAlgorithm

name := "EncryptionAtRestAlgorithm"

applicable {
    input.atRestEncryption[_]
}

compliant {
	enc := input.atRestEncryption
    data.operator == "=="
	enc.algorithm == data.target_value
}
