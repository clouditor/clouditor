package csaf

import (
	"fmt"
	"path/filepath"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gocsaf/csaf/v3/csaf"
	csafutil "github.com/gocsaf/csaf/v3/util"
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

// handleProvider tries to convert a [csaf.LoadedProviderMetadata] into an
// [ontology.SecurityAdvisoryService] as well as associated resources, such as
// its metadata document and signing keys.
func (d *csafDiscovery) handleProvider(lpmd *csaf.LoadedProviderMetadata) (resources []ontology.IsResource, err error) {
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

	// TODO(oxisto): find a sensible ID instead of this one
	serviceId := lpmd.URL + "/service"

	// Discover keys by looping through the keys provided in the PMD
	keys, keyring := d.discoverKeys(pmd.PGPKeys, serviceId)

	// Discover advisory documents from this provider
	securityAdvisoryDocuments, err := d.discoverSecurityAdvisories(lpmd, keyring, serviceId)
	if err != nil {
		return nil, fmt.Errorf("could not discover security advisories: %w", err)
	}

	var provider = &ontology.SecurityAdvisoryService{
		Id:                         serviceId,
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

	resources = append(resources, serviceMetadata, provider)
	resources = append(resources, securityAdvisoryDocuments...)
	resources = append(resources, keys...)
	return
}

func getIDsOf(documents []ontology.IsResource) (ids []string) {
	for _, d := range documents {
		ids = append(ids, d.GetId())
	}
	return
}
