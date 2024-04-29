package csaf

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/constants"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
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

// documentChecksums returns a list of [ontology.DocumentChecksum] objects, that are filled with the information
// retrieved from the checksum of the advisory file. The current implementation checks for SHA- 256 and SHA-512
// checksums located at the URL specified in the [csaf.AdvisoryFile] interface.
//
// It uses the [documentChecksum] method to retrieve the respective checksum.
func (d *csafDiscovery) documentChecksums(file csaf.AdvisoryFile, body []byte) (checksums []*ontology.DocumentChecksum) {
	var toHash = map[string]struct {
		h   hash.Hash
		alg string
	}{
		file.SHA256URL(): {h: sha256.New(), alg: constants.SHA_256},
		file.SHA512URL(): {h: sha512.New(), alg: constants.SHA_512},
	}

	for url, info := range toHash {
		checksum := d.documentChecksum(url, filepath.Base(file.URL()), body, info.alg, info.h)
		if checksum != nil {
			checksums = append(checksums, checksum)
		}
	}

	return
}

// documentChecksum retrieves a hash-based document checksum from the URL specified in checksumURL. This is usually
// constructed based on the filename by appending an extension based on the hashing algorithm, e.g. `.sha256`. It uses
// the configured client in the [csafDiscovery] struct to retrieve the checksum contents and compares it to the body, by
// invoking the specified hashing algorithm on its own.
//
// The validation result (and any resulting) errors are stored in an [ontology.DocumentChecksum] object. It will return
// nil, if the document does not exist (because then no checksum exists) or if any other (network) error occurred before
// we can retrieve the body.
func (d *csafDiscovery) documentChecksum(checksumURL, filename string, body []byte, algorithm string, h hash.Hash) *ontology.DocumentChecksum {
	// Retrieve hash from URL
	res, err := d.client.Get(checksumURL)
	if err != nil {
		// Some network error or something else happened. In this case, we cannot really produce any evidence at all
		return nil
	}

	// If the URL does not exist, we do not have a hash. In any other cases, we have some kind of response that we will
	// match against the body. From this point on, we convert all errors into ontology (validation) errors.
	if res.StatusCode == 404 {
		return nil
	}

	// Fetch the body
	hashBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &ontology.DocumentChecksum{
			Errors:    fromError(err),
			Algorithm: algorithm,
		}
	}

	// Split it into hash and filename
	hash, _, found := strings.Cut(string(hashBody), filename)
	if !found || filename == "" {
		return &ontology.DocumentChecksum{
			Errors:    fromError(errors.New("checksum file does not contain correct filename")),
			Algorithm: algorithm,
		}
	}

	// Remove any white-spaces on the "right"
	hash = strings.TrimRight(hash, " ")

	// Compare the hash
	_, _ = h.Write(body)
	want := hex.EncodeToString(h.Sum(nil)[:])
	if subtle.ConstantTimeCompare([]byte(hash), []byte(want)) == 0 {
		return &ontology.DocumentChecksum{
			Errors:    fromError(errors.New("checksum mismatch")),
			Algorithm: algorithm,
		}
	}

	// If we arrived here, everything is good
	return &ontology.DocumentChecksum{
		Errors:    nil,
		Algorithm: algorithm,
	}
}

// documentPGPSignature retrieves a signature file from the signURL and validates, whether the signatures matches the
// body based on the public keys in the keyring.
//
// The validation result (and any resulting) errors are stored in an [ontology.DocumentSignature] object. It will return
// nil, if the document does not exist (because then no signature exists) or if any other (network) error occurred
// before we can retrieve the body.
func (d *csafDiscovery) documentPGPSignature(signURL string, body []byte, keyring openpgp.EntityList) (sig *ontology.DocumentSignature) {
	var (
		res *http.Response
	)

	// Retrieve hash from URL
	res, err := d.client.Get(signURL)
	if err != nil {
		// Some network error or something else happened. In this case, we cannot really produce any evidence at all
		return nil
	}

	// If the URL does not exist, we do not have a hash. In any other cases, we have some kind of response that we will
	// match against the body. From this point on, we convert all errors into ontology (validation) errors.
	if res.StatusCode == 404 {
		return nil
	}

	// Fetch the signature (in res.Body) and use it to verify the body
	defer res.Body.Close()
	_, err = openpgp.CheckArmoredDetachedSignature(keyring, bytes.NewReader(body), res.Body, nil)
	if err != nil {
		return &ontology.DocumentSignature{
			Errors:    fromError(err),
			Algorithm: "PGP",
		}
	}

	return &ontology.DocumentSignature{
		Errors:    nil,
		Algorithm: "PGP",
	}
}

// fromError is small helper function that returns a list of [ontology.Error] objects, that can be used in various
// ontology structs based on a Go error. If the go error supports wrapping multiple errors, then the individual errors
// are added.
func fromError(err error) (errors []*ontology.Error) {
	type MultiWrapError interface {
		Unwrap() []error
	}

	if me, ok := err.(MultiWrapError); ok {
		errs := me.Unwrap()
		for _, err := range errs {
			errors = append(errors, &ontology.Error{Message: err.Error()})
		}
	} else {
		errors = append(errors, &ontology.Error{Message: err.Error()})
	}

	return
}
