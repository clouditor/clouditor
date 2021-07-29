package voc

type AtRestEncryption struct {
	*Confidentiality
	Keymanager	string `json:"keymanager"`
	Algorithm	string `json:"algorithm"`
	Enabled	bool `json:"enabled"`
}

