package providertest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

type ServiceHandler interface {
	handleIndexTxt(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory)
	handleChangesCsv(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory)
	handleAdvisory(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory)
}

func NewGoodIndexTxtWriter() ServiceHandler {
	return &goodIndexTxtWriter{}
}

// for future tests
//type errIndexTxtWriter struct {
//	statusCode int
//}

type goodIndexTxtWriter struct{}

func (good *goodIndexTxtWriter) handleIndexTxt(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory) {
	for _, advisory := range advisories {
		// write something, take URL from tracking ID
		_ = advisory.Document.Tracking.ID
	}
}

func (good *goodIndexTxtWriter) handleChangesCsv(w http.ResponseWriter, r *http.Request, advisories []*csaf.Advisory) {
	// TODO: must be sorted!
	for _, advisory := range advisories {
		line := fmt.Sprintf("\"%s\",\"%s\"\n", DocURL(advisory.Document), util.Deref(advisory.Document.Tracking.CurrentReleaseDate))
		// write something, take release from tracking current_release_data
		w.Write([]byte(line))
	}
}

func (good *goodIndexTxtWriter) handleAdvisory(w http.ResponseWriter, r *http.Request, advisory *csaf.Advisory) {
	b, err := json.Marshal(advisory)
	if err == nil {
		w.Write(b)
	}
}

func DocURL(doc *csaf.Document) string {
	// Need to parse the date
	t, _ := time.Parse(time.RFC3339, *doc.Tracking.InitialReleaseDate)
	return path.Join(strconv.FormatInt(int64(t.Year()), 10), strings.ToLower(string(util.Deref(doc.Tracking.ID)))+".json")
}
