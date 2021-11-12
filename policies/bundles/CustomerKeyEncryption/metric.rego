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
    # Check if keyUrl is set. It is only set in the customer key encryption case.
	enc.keyUrl
}