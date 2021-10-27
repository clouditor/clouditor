package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric EncryptionAtRestAlgorithm

name := "EncryptionAtRestAlgorithm"

enc := input.atRestEncryption

applicable {
    enc
}

compliant {
    data.operator == "=="
	enc.algorithm == data.target_value
}
