package voc

type Redundancy struct {
	*Availability
	Local bool // True when local redundancy is enabled
	Zone  bool // True when zone redundancy is enabled
	Geo   bool // True when geo redundancy is enabled
}
