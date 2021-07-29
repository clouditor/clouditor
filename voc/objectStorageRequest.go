package voc

type ObjectStorageRequest struct {
	*Functionality
	ObjectStorage	[]ResourceID `json:"objectStorage"`
	Source	string `json:"source"`
	Type	string `json:"type"`
}

