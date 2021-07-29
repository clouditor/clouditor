package voc

type HttpEndpoint struct {
	*Functionality
	Authenticity	*Authenticity `json:"authenticity"`
	TransportEncryption	*TransportEncryption `json:"transportEncryption"`
	Url	string `json:"url"`
	Method	string `json:"method"`
	Handler	string `json:"handler"`
	Path	string `json:"path"`
}

