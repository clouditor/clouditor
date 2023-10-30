package voc

import "time"

var KeyType = []string{"Key"}

type Key struct {
	*Resource
	*Confidentiality
	Enabled bool
	// Todo(all): Is time the appropriate type here?
	ActivationDate *time.Time
	ExpirationDate *time.Time
	// Todo(all): Think
	IsCustomerGenerated bool
	KeyType             string
	KeySize             int
	NumberOfUsages      int
}

func (*Key) Type() string {
	return "Key"
}
