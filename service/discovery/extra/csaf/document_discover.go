package csaf

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"

	"clouditor.io/clouditor/v2/api/ontology"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	csafutil "github.com/csaf-poc/csaf_distribution/v3/util"
)

func (d *csafDiscovery) discoverSecurityAdvisories(md *csaf.LoadedProviderMetadata) (documents []ontology.IsResource, err error) {
	baseURL, err := url.Parse(md.URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base URL: %w", err)
	}

	afp := csaf.NewAdvisoryFileProcessor(d.client, csafutil.NewPathEval(), md.Document, baseURL)
	err = afp.Process(func(label csaf.TLPLabel, files []csaf.AdvisoryFile) error {
		for _, f := range files {
			doc, err := d.handleAdvisory(label, f)
			if err != nil {
				return err
			}

			documents = append(documents, doc)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not process advisory files: %w", err)
	}

	return
}

func (d *csafDiscovery) handleAdvisory(label csaf.TLPLabel, file csaf.AdvisoryFile) (doc *ontology.SecurityAdvisoryDocument, err error) {
	// Create an evidence for the document
	doc = &ontology.SecurityAdvisoryDocument{
		Filetype: "JSON",
		Id:       file.URL(),
		Labels: map[string]string{
			"tlp": string(label),
		},
		Name: filepath.Base(file.URL()),
		Path: file.URL(),
	}

	// Next, we actually need to retrieve the document to check its validity
	res, err := d.client.Get(file.URL())
	if err != nil {
		// TODO: actually still need to produce an evidence that the http request was not good
		return nil, err
	}

	var raw any

	json.NewDecoder(res.Body).Decode(&raw)

	// TODO(oxisto): Check for the hashes
	msg, err := csaf.ValidateCSAF(raw)
	doc.SchemaValidation = &ontology.SchemaValidation{
		SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/csaf_json_schema.json",
		Format:    "Common Security Advisory Framework",
		Errors:    documentValidationErrors(msg),
	}

	return
}
