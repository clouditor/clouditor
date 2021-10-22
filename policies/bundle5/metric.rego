package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric EncryptionAtRestEnabled

name := "EncryptionAtRestEnabled"
metricID := 5

enc := input.atRestEncryption

applicable {
    enc
}

compliant {
    data.operator == "=="
	enc.enabled == data.target_value
}
