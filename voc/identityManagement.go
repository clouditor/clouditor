package voc

type IdentityManagement struct {
	*CloudResource
	Authorization	*Authorization `json:"authorization"`
	Authenticity	*Authenticity `json:"authenticity"`
}

