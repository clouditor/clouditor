package csaf

import (
	"fmt"
	"path/filepath"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	csafutil "github.com/csaf-poc/csaf_distribution/v3/util"
)

func (d *csafDiscovery) discoverProviders() (providers []ontology.IsResource, err error) {
	loader := csaf.NewProviderMetadataLoader(d.client)
	lpmd := loader.Load(d.domain)

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
		Path:     lpmd.URL,
		SchemaValidation: &ontology.SchemaValidation{
			Format:    "CSAF provider metadata",
			SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/provider_json_schema.json",
			Errors:    providerValidationErrors(lpmd.Messages),
		},
	}

	// TODO: Add discovery of advisory documents

	var provider = &ontology.SecurityAdvisoryService{
		Id:                          lpmd.URL + "service",
		InternetAccessibleEndpoint:  true,
		Name:                        util.Deref(pmd.Publisher.Name),
		SecurityAdvisoryDocumentIds: []string{}, // TODO: When advisory documents added, get IDs from them
		ServiceMetadataDocumentId:   util.Ref(serviceMetadata.Id),
		TransportEncryption:         d.transportEncryption(lpmd.URL),
	}

	providers = append(providers, serviceMetadata, provider)

	return
}
