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

	// Discover advisory documents from this provider
	securityAdvisoryDocuments, err := d.discoverSecurityAdvisories(lpmd)
	if err != nil {
		return nil, fmt.Errorf("could not discover security advisories: %w", err)
	}

	var provider = &ontology.SecurityAdvisoryService{
		Id:                          lpmd.URL + "service",
		InternetAccessibleEndpoint:  true,
		Name:                        util.Deref(pmd.Publisher.Name),
		SecurityAdvisoryDocumentIds: getIDsOf(securityAdvisoryDocuments),
		ServiceMetadataDocumentId:   util.Ref(serviceMetadata.Id),
		TransportEncryption:         d.transportEncryption(lpmd.URL),
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
