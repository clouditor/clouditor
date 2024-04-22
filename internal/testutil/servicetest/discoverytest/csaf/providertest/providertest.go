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
	writeJSON(w, p.pmd)
}

func (p *TrustedProvider) Domain() string {
	return p.Listener.Addr().String()
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
