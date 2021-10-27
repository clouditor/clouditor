package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric CustomerKeyEncryption

name := "CustomerKeyEncryption"

enc := input.atRestEncryption

applicable {
    enc
}

compliant {
    # Just check if keyUrl is set.
	enc.keyUrl
}