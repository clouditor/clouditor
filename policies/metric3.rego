package clouditor

default compliant = false

compliant {
	enc := input.atRestEncryption
	enc.enabled == true
}
