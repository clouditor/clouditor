package voc

type NetworkInterface struct {
	*Networking
	NetworkService	[]ResourceID `json:"networkService"`
	AccessRestriction	*AccessRestriction `json:"accessRestriction"`
}

