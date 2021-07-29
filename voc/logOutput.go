package voc

type LogOutput struct {
	*Functionality
	Logging	[]ResourceID `json:"logging"`
	Call	string `json:"call"`
	Value	string `json:"value"`
}

