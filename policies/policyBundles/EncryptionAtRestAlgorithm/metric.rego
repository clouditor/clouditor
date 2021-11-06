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
    compare(data.operator, data.target_value, enc.algorithm)
}
