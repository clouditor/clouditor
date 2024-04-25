package csaf

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	// Next, we actually need to retrieve the document to check its validity
	res, err := d.client.Get(file.URL())
	if err != nil {
		// TODO: actually still need to produce an evidence that the http request was not good. This goes for all errors I guess?
		return nil, err
	}

	var raw any
	var advisory csaf.Advisory

	err = json.NewDecoder(res.Body).Decode(&raw)
	if err != nil {
		return nil, err
	}

	// TODO(oxisto): Check for the hashes
	msg, err := csaf.ValidateCSAF(raw)
	if err != nil {
		return nil, err
	}

	// ReMarshal into a struct that we can actually work with
	err = csafutil.ReMarshalJSON(&advisory, raw)
	if err != nil {
		return nil, err
	}

	time, err := time.Parse(time.RFC3339, util.Deref(advisory.Document.Tracking.InitialReleaseDate))
	if err != nil {
		return nil, err
	}

	// Create an evidence for the document
	doc = &ontology.SecurityAdvisoryDocument{
		Filetype: "JSON",
		Id:       string(*advisory.Document.Tracking.ID),
		Labels: map[string]string{
			"tlp": string(label),
		},
		Name: util.Deref(advisory.Document.Title),
		DocumentLocation: &ontology.DocumentLocation{
			Type: &ontology.DocumentLocation_RemoteDocumentLocation{
				RemoteDocumentLocation: &ontology.RemoteDocumentLocation{
					Path:                file.URL(),
					TransportEncryption: transportEncryption(res.TLS),
					Authenticity:        clientAuthenticity(res),
				},
			},
		},
		CreationTime: timestamppb.New(time),
		SchemaValidation: &ontology.SchemaValidation{
			SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/csaf_json_schema.json",
			Format:    "Common Security Advisory Framework",
			Errors:    documentValidationErrors(msg),
		},
	}

	return
}
