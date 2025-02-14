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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"

	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/gocsaf/csaf/v3/csaf"
	csafutil "github.com/gocsaf/csaf/v3/util"
)

type TrustedProvider struct {
	*httptest.Server
	PMD   *csaf.ProviderMetadata
	idxw  ServiceHandler
	feeds map[csaf.TLPLabel][]*csaf.Advisory

	Keyring openpgp.EntityList
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

	// Create a new OpenPGP key pair
	key, _ := openpgp.NewEntity("test", "test", "test", nil)
	p.Keyring = append(p.Keyring, key)

	mux.HandleFunc("/.well-known/csaf/provider-metadata.json", p.handlePMD)

	p.PMD = csaf.NewProviderMetadataDomain(fmt.Sprintf("https://%s", p.Domain()), nil)

	// We need to provide one index.txt per feed. So far, we only support index.txt, no ROLIE
	for feed := range p.feeds {
		feedURL := fmt.Sprintf("/.well-known/csaf/%s/", strings.ToLower(string(feed)))
		mux.HandleFunc(feedURL, p.handleFeed)
		p.PMD.Distributions = append(p.PMD.Distributions, csaf.Distribution{
			DirectoryURL: fmt.Sprintf("https://%s/%s", p.Domain(), feedURL),
			Rolie:        nil,
		})
	}

	for _, key := range p.Keyring {
		fp := hex.EncodeToString(key.PrimaryKey.Fingerprint)
		p.PMD.PGPKeys = append(p.PMD.PGPKeys, csaf.PGPKey{
			Fingerprint: csaf.Fingerprint(fp),
			URL:         util.Ref("https://" + p.Domain() + "/.well-known/csaf/opengpg/" + fp + ".asc"),
		})
		mux.HandleFunc("/.well-known/csaf/opengpg/"+fp+".asc", p.handleKey)
	}

	// Apply PMD functions
	for _, fn := range fns {
		fn(p.PMD)
	}

	return
}

func (p *TrustedProvider) handlePMD(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, p.PMD)
}

func (p *TrustedProvider) handleKey(w http.ResponseWriter, r *http.Request) {
	// Find key from fingerprint
	file := filepath.Base(r.URL.Path)
	fingerprint := strings.TrimSuffix(file, filepath.Ext(file))

	idx := slices.IndexFunc(p.Keyring, func(key *openpgp.Entity) bool {
		return hex.EncodeToString(key.PrimaryKey.Fingerprint) == fingerprint
	})
	if idx != -1 {
		s, err := openpgp.WriteArmoredKey(p.Keyring[idx])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte(s))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
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
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file := filepath.Base(r.URL.Path)

	if file == "index.txt" {
		p.idxw.handleIndexTxt(w, r, advisories, p)
	} else if file == "changes.csv" {
		p.idxw.handleChangesCsv(w, r, advisories, p)
	} else {
		ext := filepath.Ext(file)
		var id string
		var handler func(http.ResponseWriter, *http.Request, *csaf.Advisory, *TrustedProvider)

		if ext == ".sha256" {
			id = strings.TrimSuffix(file, ".json"+filepath.Ext(file))
			handler = p.idxw.handleSHA256
		} else if ext == ".sha512" {
			id = strings.TrimSuffix(file, ".json"+filepath.Ext(file))
			handler = p.idxw.handleSHA512
		} else if ext == ".asc" {
			id = strings.TrimSuffix(file, ".json"+filepath.Ext(file))
			handler = p.idxw.handleSignature
		} else {
			id = strings.TrimSuffix(file, filepath.Ext(file))
			handler = p.idxw.handleAdvisory
		}

		// Find the advisory
		idx := slices.IndexFunc(advisories, func(doc *csaf.Advisory) bool {
			return strings.ToLower(string(util.Deref(doc.Document.Tracking.ID))) == id
		})
		if idx == -1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler(w, r, advisories[idx], p)
	}
}

func (p *TrustedProvider) Domain() string {
	return p.Listener.Addr().String()
}

// WellKnownProviderURL returns the URL of the provider-metadata.json at the well-known address.
func (p *TrustedProvider) WellKnownProviderURL() string {
	return "https://" + p.Domain() + "/.well-known/csaf/provider-metadata.json"
}

// DocumentAny returns the [csaf.ProviderMetadata] as "any" (actually a map[string]interface{}) so that it can be used
// in [csaf.LoadedProviderMetadata].
func (p *TrustedProvider) DocumentAny() (m map[string]interface{}) {
	err := csafutil.ReMarshalJSON(&m, p.PMD)
	if err != nil {
		panic(err)
	}
	return m
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
