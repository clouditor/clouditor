package voc

type HttpRequest struct {
	*Functionality
	HttpEndpoint	*HttpEndpoint `json:"httpEndpoint"`
	Call	string `json:"call"`
}

