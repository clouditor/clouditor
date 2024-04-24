// Copyright 2024 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package providertest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"

	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

type TrustedProvider struct {
	*httptest.Server
	pmd   *csaf.ProviderMetadata
	idxw  ServiceHandler
	feeds map[csaf.TLPLabel][]*csaf.Advisory
}

// TODO: convert into functional-style options?
func NewTrustedProvider(
	feeds map[csaf.TLPLabel][]*csaf.Advisory,
	idxw ServiceHandler,
	fns ...func(*csaf.ProviderMetadata),
) (p *TrustedProvider) {
	mux := http.NewServeMux()

	p = &TrustedProvider{
		idxw:  idxw,
		feeds: feeds,
	}
	p.Server = httptest.NewTLSServer(mux)
	p.Server.EnableHTTP2 = true

	mux.HandleFunc("/.well-known/csaf/provider-metadata.json", p.handlePMD)

	p.pmd = csaf.NewProviderMetadataDomain(fmt.Sprintf("https://%s", p.Domain()), nil)

	// We need to provide one index.txt per feed. So far, we only support index.txt, no ROLIE
	for feed := range p.feeds {
		feedURL := fmt.Sprintf("/.well-known/csaf/%s/", strings.ToLower(string(feed)))
		mux.HandleFunc(feedURL, p.handleFeed)
		p.pmd.Distributions = append(p.pmd.Distributions, csaf.Distribution{
			DirectoryURL: fmt.Sprintf("https://%s/%s", p.Domain(), feedURL),
			Rolie:        nil,
		})
	}

	// Apply PMD functions
	for _, fn := range fns {
		fn(p.pmd)
	}

	return
}

func (p *TrustedProvider) handlePMD(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, p.pmd)
}

func (p *TrustedProvider) handleFeed(w http.ResponseWriter, r *http.Request) {
	_, feed, ok := strings.Cut(filepath.Dir(r.URL.Path), "/.well-known/csaf/")
	if !ok {
		return
	}

	// make sure we don't end up with <feed>/<year>
	before, _, ok := strings.Cut(feed, "/")
	if ok {
		feed = before
	}

	advisories := p.advisoriesFor(feed)
	if advisories == nil {
		w.WriteHeader(404)
		return
	}

	file := filepath.Base(r.URL.Path)

	if file == "index.txt" {
		p.idxw.handleIndexTxt(w, r, advisories)
	} else if file == "changes.csv" {
		p.idxw.handleChangesCsv(w, r, advisories)
	} else {
		idx := slices.IndexFunc(advisories, func(doc *csaf.Advisory) bool {
			return strings.ToLower(string(util.Deref(doc.Document.Tracking.ID))) == strings.TrimSuffix(file, filepath.Ext(file))
		})
		if idx == -1 {
			w.WriteHeader(404)
			return
		}
		p.idxw.handleAdvisory(w, r, advisories[idx])
	}
}

func (p *TrustedProvider) Domain() string {
	return p.Listener.Addr().String()
}

func (p *TrustedProvider) advisoriesFor(feed string) []*csaf.Advisory {
	return p.feeds[csaf.TLPLabel(strings.ToUpper(feed))]
}

func writeJSON(w http.ResponseWriter, msg any) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
