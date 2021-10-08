package clouditor

default compliant = false

# this is an implementation of metric EncryptionAtRestAlgorithm


compliant {
	enc := input.atRestEncryption
	goodAlgorithm(enc)
}

goodAlgorithm(enc) {
	enc.algorithm == "AES-128"
}

goodAlgorithm(tls) {
	enc.algorithm == "AES-256"
}
