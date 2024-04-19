package providertest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

type TrustedProvider struct {
	*httptest.Server
	pmd *csaf.ProviderMetadata
}

func NewTrustedProvider(fns ...func(*csaf.ProviderMetadata)) (p *TrustedProvider) {
	mux := http.NewServeMux()

	p = &TrustedProvider{}
	p.Server = httptest.NewTLSServer(mux)
	p.Server.EnableHTTP2 = true

	mux.HandleFunc("/.well-known/csaf/provider-metadata.json", p.handlePMD)
	p.pmd = csaf.NewProviderMetadataDomain(fmt.Sprintf("https://%s", p.Domain()), []csaf.TLPLabel{csaf.TLPLabelWhite})

	// Apply PMD functions
	for _, fn := range fns {
		fn(p.pmd)
	}

	return
}

func (p *TrustedProvider) handlePMD(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(p.pmd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (p *TrustedProvider) Domain() string {
	return p.Listener.Addr().String()
}
