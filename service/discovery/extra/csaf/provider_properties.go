package csaf

import (
	"clouditor.io/clouditor/v2/api/ontology"

	"github.com/gocsaf/csaf/v3/csaf"
)

func (d *csafDiscovery) providerTransportEncryption(url string) *ontology.TransportEncryption {
	res, err := d.client.Get(url)
	if err != nil {
		return &ontology.TransportEncryption{
			Enabled: new(false),
		}
	}

	return transportEncryption(res.TLS)
}

func providerValidationErrors(messages csaf.ProviderMetadataLoadMessages) (errs []*ontology.Error) {
	for _, m := range messages {
		errs = append(errs, &ontology.Error{Message: new(m.Message)})
	}
	return
}
