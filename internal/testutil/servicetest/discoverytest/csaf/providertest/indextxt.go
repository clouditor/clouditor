package providertest

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gocsaf/csaf/v3/csaf"
	log "github.com/sirupsen/logrus"
)

type ServiceHandler interface {
	handleIndexTxt(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory, p *TrustedProvider)
	handleChangesCsv(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory, p *TrustedProvider)
	handleAdvisory(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory, p *TrustedProvider)
	handleSHA256(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory, p *TrustedProvider)
	handleSHA512(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory, p *TrustedProvider)
	handleSignature(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory, p *TrustedProvider)
}

func NewGoodIndexTxtWriter() ServiceHandler {
	return &goodIndexTxtWriter{}
}

// for future tests
//type errIndexTxtWriter struct {
//	statusCode int
//}

type goodIndexTxtWriter struct{}

func (good *goodIndexTxtWriter) handleIndexTxt(_ http.ResponseWriter, _ *http.Request, advisories []*csaf.Advisory, _ *TrustedProvider) {
	for _, advisory := range advisories {
		// write something, take URL from tracking ID
		_ = advisory.Document.Tracking.ID
	}
}

func (good *goodIndexTxtWriter) handleChangesCsv(w http.ResponseWriter, _ *http.Request, advisories []*csaf.Advisory, _ *TrustedProvider) {
	// TODO: must be sorted!
	for _, advisory := range advisories {
		line := fmt.Sprintf("\"%s\",\"%s\"\n", DocURL(advisory.Document), util.Deref(advisory.Document.Tracking.CurrentReleaseDate))
		// write something, take release from tracking current_release_data
		_, err := w.Write([]byte(line))
		// Maybe do better error handling
		if err != nil {
			log.Warnf("Could not write csv: %s", err.Error())
		}
	}
}

func (good *goodIndexTxtWriter) handleAdvisory(w http.ResponseWriter, _ *http.Request, advisory *csaf.Advisory, _ *TrustedProvider) {
	b, err := json.Marshal(advisory)
	if err == nil {
		_, err = w.Write(b)
		// Maybe do better error handling
		if err != nil {
			log.Warnf("Could not write: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (good *goodIndexTxtWriter) handleSHA256(w http.ResponseWriter, _ *http.Request, advisory *csaf.Advisory, _ *TrustedProvider) {
	good.handleHash(w, advisory, sha256.New())
}

func (good *goodIndexTxtWriter) handleSHA512(w http.ResponseWriter, _ *http.Request, advisory *csaf.Advisory, _ *TrustedProvider) {
	good.handleHash(w, advisory, sha512.New())
}

func (good *goodIndexTxtWriter) handleHash(w http.ResponseWriter, advisory *csaf.Advisory, h hash.Hash) {
	var (
		err      error
		body     []byte
		checksum []byte
	)

	// Retrieve the body
	body, _ = json.Marshal(advisory)
	_, _ = h.Write(body)
	checksum = h.Sum(nil)

	_, err = w.Write([]byte(fmt.Sprintf("%s %s",
		hex.EncodeToString(checksum),
		strings.ToLower(string(util.Deref(advisory.Document.Tracking.ID)))+".json")),
	)
	if err != nil {
		log.Warnf("Could not write: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (good *goodIndexTxtWriter) handleSignature(w http.ResponseWriter, _ *http.Request, advisory *csaf.Advisory, p *TrustedProvider) {
	var (
		err  error
		body []byte
	)

	// Retrieve the body
	body, _ = json.Marshal(advisory)

	// Sign it
	err = openpgp.ArmoredDetachSignText(w, p.Keyring[0], bytes.NewReader(body), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DocURL(doc *csaf.Document) string {
	// Need to parse the date
	t, _ := time.Parse(time.RFC3339, *doc.Tracking.InitialReleaseDate)
	return path.Join(strconv.FormatInt(int64(t.Year()), 10), strings.ToLower(string(util.Deref(doc.Tracking.ID)))+".json")
}
