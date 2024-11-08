// Copyright 2021 Fraunhofer AISEC
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

package cmc

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"

	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	ar "github.com/Fraunhofer-AISEC/cmc/attestationreport"
	atls "github.com/Fraunhofer-AISEC/cmc/attestedtls"
	"github.com/Fraunhofer-AISEC/cmc/cmc"
)

const ()

// discoverReports discovers the attestation reports from the CMC
func (d *cmcDiscovery) discoverReports() ([]ontology.IsResource, error) {
	var (
		list []ontology.IsResource
	)

	// Read CA from filesystem
	// TODO(all): Should be removed in future, just for testing
	ca, err := os.ReadFile(d.capemPath)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate from path '%s': %w", d.capemPath, err)
	}
	log.Debugf("Certificate read from path: %s", d.capemPath)

	log.Debug("Initializing CMC")
	cmc, err := cmc.NewCmc(&cmc.Config{
		Api: "libapi",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create CMC config: %v", err)
	}

	// Add root CA
	log.Debug("Adding CA")
	roots := x509.NewCertPool()
	success := roots.AppendCertsFromPEM(ca)
	if !success {
		return nil, fmt.Errorf("could not add cert '%s' to root CAs", ca)
	}

	// Create TLS config with root CA only
	tlsConf := &tls.Config{
		RootCAs:       roots,
		Renegotiation: tls.RenegotiateNever,
	}

	conn, err := atls.Dial("tcp", d.cmcAddr, tlsConf,
		atls.WithCmcCa(ca),
		atls.WithCmcApi(atls.CmcApi_Lib),
		atls.WithMtls(false),
		atls.WithAttest("server"),
		atls.WithResultCb(func(result *ar.VerificationResult) {
			// TODO (anatheka): Return error
			r, err := handleReport(*result)
			if err != nil {
				log.Errorf("could not handle attestation report: %v", err)
			}

			log.Debug("attestation report: ", result)
			list = append(list, r)
		}),
		atls.WithCmc(cmc))
	if err != nil {
		return nil, fmt.Errorf("could not get attestation report: %v", err)
	}
	defer conn.Close()

	return list, nil
}

// TODO(anatheka): Maybe we should call this out of the handleVirtualMachine method.
func handleReport(result ar.VerificationResult) (*ontology.VirtualMachine, error) {
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal integrity result: %w", err)
	}

	resource := &ontology.VirtualMachine{
		Id:   result.Prover,
		Name: result.Prover,
		// CreationTime: , // TODO: TBD
		// GeoLocation: ,// TODO: TBD
		Raw: string(raw),
		RemoteAttestation: &ontology.RemoteAttestation{
			Enabled:      true,
			Status:       result.Success,
			CreationTime: util.Timestamp(result.Created),
		},
	}

	return resource, nil
}