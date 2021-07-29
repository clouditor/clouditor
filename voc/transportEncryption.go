package voc

type TransportEncryption struct {
	*Confidentiality
	Enforced	bool `json:"enforced"`
	Enabled	bool `json:"enabled"`
	TlsVersion	string `json:"tlsVersion"`
	Algorithm	string `json:"algorithm"`
}

