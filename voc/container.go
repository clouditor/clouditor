package voc

type Container struct {
	*Compute
	Image	[]ResourceID `json:"image"`
}

