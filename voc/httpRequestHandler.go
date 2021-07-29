package voc

type HttpRequestHandler struct {
	*Functionality
	Application	*Application `json:"application"`
	HttpEndpoint	*[]HttpEndpoint `json:"httpEndpoint"`
	Path	string `json:"path"`
}

