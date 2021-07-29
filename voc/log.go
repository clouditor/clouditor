package voc

type Log struct {
	*Auditing
	Activated	bool `json:"activated"`
}

