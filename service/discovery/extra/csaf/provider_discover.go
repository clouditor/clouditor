package csaf

import (
	"fmt"
	"path/filepath"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	csafutil "github.com/csaf-poc/csaf_distribution/v3/util"
)

func (d *csafDiscovery) discoverProviders() (providers []ontology.IsResource, err error) {
	loader := csaf.NewProviderMetadataLoader(d.client)
	lpmds := loader.Enumerate(d.domain)

	for _, lpmd := range lpmds {
		// Handle the single PMD files that were discovered It can happen that the PMD from
		// the well-known URL and the first one defined in the security.txt are the same,
		// so the evidence would be created two times
		res, err := d.handleProvider(lpmd)
		if err != nil {
			return nil, fmt.Errorf("could not discover security provider: %w", err)
		}
		// Add all discovered resources to the providers
		providers = append(providers, res...)
	}

	return
}

func (d *csafDiscovery) handleProvider(lpmd *csaf.LoadedProviderMetadata) (providers []ontology.IsResource, err error) {
	if !lpmd.Valid() {
		// TODO(oxisto): Even if the PMD is invalid, we still need to create an evidence for it!
		return nil, fmt.Errorf("could not load provider-metadata.json from %s", d.domain)
	}

	// Convert it to a csaf.ProviderMetadata struct for simpler access
	var pmd csaf.ProviderMetadata

	err = csafutil.ReMarshalJSON(&pmd, lpmd.Document)
	if err != nil {
		err = fmt.Errorf("could not convert provider metadata to struct: %w", err)
		return
	}

	serviceMetadata := &ontology.ServiceMetadataDocument{
		Filetype: "JSON",
		Id:       lpmd.URL,
		Name:     filepath.Base(lpmd.URL),
		DocumentLocation: &ontology.DocumentLocation{
			Type: &ontology.DocumentLocation_RemoteDocumentLocation{
				RemoteDocumentLocation: &ontology.RemoteDocumentLocation{
					Path:                lpmd.URL,
					TransportEncryption: d.providerTransportEncryption(lpmd.URL),
				},
			},
		},
		SchemaValidation: &ontology.SchemaValidation{
			Format:    "CSAF provider metadata",
			SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/provider_json_schema.json",
			Errors:    providerValidationErrors(lpmd.Messages),
		},
		Raw: discovery.Raw(pmd),
	}

	keys := d.discoverKeys(pmd.PGPKeys)
	providers = append(providers, keys...)

	keyring := openpgp.EntityList{}

	for _, keyinfo := range pmd.PGPKeys {
		key, err := d.fetchKey(keyinfo)
		if err != nil {
			return nil, err
		}
		keyring = append(keyring, key)
	}

	// Discover advisory documents from this provider
	securityAdvisoryDocuments, err := d.discoverSecurityAdvisories(lpmd, keyring)
	if err != nil {
		return nil, fmt.Errorf("could not discover security advisories: %w", err)
	}

	var provider = &ontology.SecurityAdvisoryService{
		Id:                         lpmd.URL + "service",
		InternetAccessibleEndpoint: true,
		Name:                       util.Deref(pmd.Publisher.Name),
		// TODO: actually put document in correct feed
		SecurityAdvisoryFeeds: []*ontology.SecurityAdvisoryFeed{
			{
				SecurityAdvisoryDocumentIds: getIDsOf(securityAdvisoryDocuments),
			},
		},
		ServiceMetadataDocumentId: util.Ref(serviceMetadata.Id),
		TransportEncryption:       serviceMetadata.DocumentLocation.GetRemoteDocumentLocation().GetTransportEncryption(),
		KeyIds:                    getIDsOf(keys),
		Raw:                       discovery.Raw(lpmd),
	}

	providers = append(providers, serviceMetadata, provider)
	providers = append(providers, securityAdvisoryDocuments...)
	return
}

func getIDsOf(documents []ontology.IsResource) (ids []string) {
	for _, d := range documents {
		ids = append(ids, d.GetId())
	}
	return
}
