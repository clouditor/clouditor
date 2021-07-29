package voc

type OTPBasedAuthentication struct {
	*Authenticity
	Activated	bool `json:"activated"`
}

