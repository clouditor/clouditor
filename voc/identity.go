package voc

type Identity struct {
	*IdentityManagement
	PasswordBasedAuthentication	*PasswordBasedAuthentication `json:"passwordBasedAuthentication"`
	OTPBasedAuthentication	*OTPBasedAuthentication `json:"oTPBasedAuthentication"`
}

