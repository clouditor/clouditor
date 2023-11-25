package voc

var KeyType = []string{"Key"}

type Key struct {
	*Resource
	*Confidentiality
	Enabled bool `json:"enabled"`
	// Todo(all): Is time the appropriate type here?
	ActivationDate int64  `json:"activationDate"`
	ExpirationDate int64  `json:"expirationDate"`
	KeyType        string `json:"keyType"`
	KeySize        int    `json:"keySize"`
	NumberOfUsages int    `json:"numberOfUsages"`
}

func (*Key) Type() string {
	return "Key"
}
