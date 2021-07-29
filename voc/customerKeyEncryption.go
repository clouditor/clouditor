package voc

type CustomerKeyEncryption struct {
	*AtRestEncryption
	KeyUrl	string `json:"keyUrl"`
}

