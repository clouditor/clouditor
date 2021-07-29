package voc

type Storage struct {
	*CloudResource
	AtRestEncryption	*AtRestEncryption `json:"atRestEncryption"`
}

