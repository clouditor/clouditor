package voc

type HttpEndpoint struct {
	rdfType struct{} `quad:"@type > cloud:HttpEndpoint"`

	URL string `json:"@id"`

	// somehow will not be persisted and is confused with AtRestEncryption
	TransportEncryption *TransportEncryption `quad:"cloud:transportEncryption"`
}
