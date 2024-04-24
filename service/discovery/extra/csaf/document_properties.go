package csaf

import (
	"crypto/tls"

	"clouditor.io/clouditor/v2/api/ontology"
)

func documentValidationErrors(messages []string) (errs []*ontology.Error) {
	for _, m := range messages {
		errs = append(errs, &ontology.Error{Message: m})
	}
	return
}

func transportEncryption(state *tls.ConnectionState) (te *ontology.TransportEncryption) {
	te = &ontology.TransportEncryption{
		Enabled: state != nil,
	}

	if state != nil {
		if state.Version == tls.VersionTLS10 {
			te.ProtocolVersion = 1.0
		} else if state.Version == tls.VersionTLS11 {
			te.ProtocolVersion = 1.1
		} else if state.Version == tls.VersionTLS12 {
			te.ProtocolVersion = 1.2
		} else if state.Version == tls.VersionTLS13 {
			te.ProtocolVersion = 1.3
		}
	}

	return te
}
