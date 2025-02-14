package csaf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gocsaf/csaf/v3/csaf"
	csafutil "github.com/gocsaf/csaf/v3/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (d *csafDiscovery) discoverSecurityAdvisories(md *csaf.LoadedProviderMetadata, keyring openpgp.EntityList, parentId string) (documents []ontology.IsResource, err error) {
	baseURL, err := url.Parse(md.URL)
	if err != nil {
		return nil, fmt.Errorf("could not parse base URL: %w", err)
	}

	afp := csaf.NewAdvisoryFileProcessor(d.client, csafutil.NewPathEval(), md.Document, baseURL)
	err = afp.Process(func(label csaf.TLPLabel, files []csaf.AdvisoryFile) error {
		for _, f := range files {
			doc, err := d.handleAdvisory(label, f, keyring, parentId)
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

func (d *csafDiscovery) handleAdvisory(label csaf.TLPLabel, file csaf.AdvisoryFile, keyring openpgp.EntityList, parentId string) (doc *ontology.SecurityAdvisoryDocument, err error) {
	// Next, we actually need to retrieve the document to check its validity
	res, err := d.client.Get(file.URL())
	if err != nil {
		// TODO: actually still need to produce an evidence that the http request was not good. This goes for all errors I guess?
		return nil, err
	}

	var (
		raw      any
		advisory csaf.Advisory
		body     []byte
	)

	body, err = io.ReadAll(res.Body)
	if err != nil {
		// TODO: add to validation error?
		return nil, err
	}

	err = json.Unmarshal(body, &raw)
	if err != nil {
		// TODO: add to validation error?
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

	t, err := time.Parse(time.RFC3339, util.Deref(advisory.Document.Tracking.InitialReleaseDate))
	if err != nil {
		return nil, err
	}

	// Create an evidence for the document
	doc = &ontology.SecurityAdvisoryDocument{
		Filetype: "JSON",
		Id:       string(util.Deref(advisory.Document.Tracking.ID)),
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
		CreationTime: timestamppb.New(t),
		SchemaValidation: &ontology.SchemaValidation{
			SchemaUrl: "https://docs.oasis-open.org/csaf/csaf/v2.0/csaf_json_schema.json",
			Format:    "Common Security Advisory Framework",
			Errors:    documentValidationErrors(msg),
		},
		DocumentChecksums: d.documentChecksums(file, body),
		DocumentSignatures: []*ontology.DocumentSignature{
			d.documentPGPSignature(file.SignURL(), body, keyring),
		},
		Raw:      discovery.Raw(doc),
		ParentId: &parentId,
	}

	return
}
