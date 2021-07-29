package voc

type LoadBalancer struct {
	*NetworkService
	AccessRestriction *AccessRestriction `json:"accessRestriction"`
	HttpEndpoint      *[]HttpEndpoint    `json:"httpEndpoint"`
	Url               string             `json:"url"`
}
