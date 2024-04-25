package csaf

import (
	"crypto/tls"
	"net/http"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
)

func documentValidationErrors(messages []string) (errs []*ontology.Error) {
	for _, m := range messages {
		errs = append(errs, &ontology.Error{Message: m})
	}
	return
}

// transportEncryption extracts the properties needed for a [ontology.TransportEncryption] out of a
// [tls.ConnectionState].
func transportEncryption(state *tls.ConnectionState) (te *ontology.TransportEncryption) {
	te = &ontology.TransportEncryption{}

	if state != nil {
		te.Enabled = true
		if state.Version == tls.VersionTLS10 {
			te.ProtocolVersion = 1.0
		} else if state.Version == tls.VersionTLS11 {
			te.ProtocolVersion = 1.1
		} else if state.Version == tls.VersionTLS12 {
			te.ProtocolVersion = 1.2
		} else if state.Version == tls.VersionTLS13 {
			te.ProtocolVersion = 1.3
		}

		te.Protocol = constants.TLS
		cs := cipherSuite(state.CipherSuite)
		if cs != nil {
			te.CipherSuites = append(te.CipherSuites, cs)
		}
	}

	return te
}

// cipherSuite builds an [ontology.CipherSuite] object out of the cipher suite identifier of the tls package, e.g.
// [tls.TLS_AES_128_GCM_SHA256].
func cipherSuite(id uint16) *ontology.CipherSuite {
	if id == tls.TLS_AES_128_GCM_SHA256 {
		return &ontology.CipherSuite{
			SessionCipher: constants.AES_128_GCM,
			MacAlgorithm:  constants.SHA_256,
		}
	} else if id == tls.TLS_AES_256_GCM_SHA384 {
		return &ontology.CipherSuite{
			SessionCipher: constants.AES_256_GCM,
			MacAlgorithm:  constants.SHA_384,
		}
	}
	return nil
}

func clientAuthenticity(res *http.Response) *ontology.Authenticity {
	// If we did not have any authorization header on our client and the request was successful, we have
	// "NoAuthentication"
	if res.Request.Header.Get("authorization") == "" && res.StatusCode != http.StatusUnauthorized {
		return &ontology.Authenticity{
			Type: &ontology.Authenticity_NoAuthentication{},
		}
	}

	return nil
}
