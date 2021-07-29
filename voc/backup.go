package voc

type Backup struct {
	*Availability
	Activated	bool `json:"activated"`
	Policy	string `json:"policy"`
}

