package voc

type ObjectStorage struct {
	*Storage
	HttpEndpoint	*HttpEndpoint `json:"httpEndpoint"`
}

