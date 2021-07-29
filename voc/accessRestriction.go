package voc

type AccessRestriction struct {
	*Authorization
	Inbound	bool `json:"inbound"`
	RestrictedPorts	string `json:"restrictedPorts"`
}

