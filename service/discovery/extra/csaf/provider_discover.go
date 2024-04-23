package csaf

import (
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"fmt"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"net/url"
	"path/filepath"

	csaf_util "github.com/csaf-poc/csaf_distribution/v3/util"
)

func (d *csafDiscovery) discoverProviders(md *csaf.LoadedProviderMetadata) (providers []ontology.IsResource, err error) {
	pmd, ok := md.Document.(csaf.ProviderMetadata)
	if !ok {
		err = fmt.Errorf("could not assert metadata pmd: %v", err)
		return
	}
	serviceMetadata := &ontology.ServiceMetadataDocument{
		Filetype: "JSON",
		Id:       md.URL,
		Name:     filepath.Base(md.URL),
		Path:     md.URL,
		SchemaValidation: &ontology.SchemaValidation{
			Format:    "CSAF provider metadata",
			SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/provider_json_schema.json",
			Errors:    convertToErrors(md.Messages), // TODO(lebogg):convert
		},
	}

	// get getSecurityAdvisoryDocuments
	securityAdvisoryDocuments := d.getSecurityAdvisoryDocuments(md)

	var provider = &ontology.SecurityAdvisoryService{
		Id:                          md.URL + "service",
		InternetAccessibleEndpoint:  true,
		Name:                        util.Deref(pmd.Publisher.Name),
		SecurityAdvisoryDocumentIds: getIDsOf(securityAdvisoryDocuments),
		ServiceMetadataDocumentId:   util.Ref(serviceMetadata.Id),
		TransportEncryption:         d.checkTransportEncryption(md.URL),
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

func (d *csafDiscovery) getSecurityAdvisoryDocuments(md *csaf.LoadedProviderMetadata) (documents []ontology.IsResource) {
	baseURL, err := url.Parse(md.URL)
	if err != nil {
		return nil // TODO(lebogg): handle
	}

	afp := csaf.NewAdvisoryFileProcessor(d.client, csaf_util.NewPathEval(), md.Document, baseURL)
	err = afp.Process(func(label csaf.TLPLabel, files []csaf.AdvisoryFile) error {
		for _, f := range files {
			documents = append(documents, &ontology.SecurityAdvisoryDocument{
				Filetype: "JSON",
				Id:       f.URL(),
				Labels: map[string]string{
					"tlp": string(label),
				},
				Name: filepath.Base(f.URL()),
				Path: f.URL(),
			})
		}
		return nil
	})
	if err != nil {
		log.Errorf("Could not process advisory files: %v", err)
	}

	return
}

func convertToErrors(messages csaf.ProviderMetadataLoadMessages) (errs []*ontology.Error) {
	for _, m := range messages {
		errs = append(errs, &ontology.Error{Message: m.Message})
	}
	return
}

func (d *csafDiscovery) checkTransportEncryption(url string) *ontology.TransportEncryption {
	res, err := d.client.Get(url)
	if err != nil {
		return &ontology.TransportEncryption{
			Enabled: false,
		}
	}
	return &ontology.TransportEncryption{
		Enabled: res.TLS != nil,
	}
}
