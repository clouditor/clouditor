package voc

type GeoLocation struct {
	*Availability
	Region	string `json:"region"`
}

