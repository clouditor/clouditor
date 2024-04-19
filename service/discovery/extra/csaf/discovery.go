// csaf contains a discover that discovery security advisory information from a CSAF trusted provider
package csaf

import (
	"fmt"
	"net/http"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "azure-discovery")
}

type csafDiscovery struct {
	domain string
	csID   string
	client *http.Client
}

type DiscoveryOption func(a *csafDiscovery)

func WithProviderDomain(domain string) DiscoveryOption {
	return func(a *csafDiscovery) {
		a.domain = domain
	}
}

func WithCloudServiceID(csID string) DiscoveryOption {
	return func(a *csafDiscovery) {
		a.csID = csID
	}
}

func NewCSAFDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &csafDiscovery{
		domain: "wid.cert-bund.de",
		client: http.DefaultClient,
	}

	// Apply options
	for _, opt := range opts {
		opt(d)
	}

	return d
}

func (*csafDiscovery) Name() string {
	return "CSAF Trusted Provider Discovery"
}

func (*csafDiscovery) Description() string {
	return "Discovery CSAF documents from a CSAF trusted provider"
}

func (d *csafDiscovery) CloudServiceID() string {
	return d.csID
}

func (d *csafDiscovery) List() (list []ontology.IsResource, err error) {
	log.Info("Fetching CSAF documents from provider")

	loader := csaf.NewProviderMetadataLoader(d.client)
	metadata := loader.Load(d.domain)
	_ = metadata

	if !metadata.Valid() {
		return nil, fmt.Errorf("could not load provider-metadata.json from %s", d.domain)
	}

	// TODO: actually discover evidences in future PR

	return nil, nil
}
