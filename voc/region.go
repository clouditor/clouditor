package voc

type HasGeoLocation interface {
	GeoLocation() GeoLocation
}

type GeoLocation struct {
}
