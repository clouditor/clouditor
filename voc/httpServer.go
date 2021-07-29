package voc

type HttpServer struct {
	*Framework
	HttpRequestHandler	*HttpRequestHandler `json:"httpRequestHandler"`
}

