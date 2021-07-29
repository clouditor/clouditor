package voc

type PasswordBasedAuthentication struct {
	*Authenticity
	Activated	bool `json:"activated"`
}

