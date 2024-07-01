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
	"time"

	"clouditor.io/clouditor/v2/api/ontology"
	"google.golang.org/protobuf/types/known/timestamppb"

	ar "github.com/Fraunhofer-AISEC/cmc/attestationreport"
	atls "github.com/Fraunhofer-AISEC/cmc/attestedtls"
	"github.com/Fraunhofer-AISEC/cmc/cmc"
)

const (
	timeoutSec = 10
	capemPath  = "local/certificate_remote_attestation.pem"
)

// discoverReports discovers the attestation reports from the CMC
func (d *cmcDiscovery) discoverReports() ([]ontology.IsResource, error) {
	var (
		list []ontology.IsResource
		ca   = []byte{}
	)

	// Read CA from filesystem
	// TODO(all): Should be removed in future, just for testing
	ca, err := os.ReadFile(capemPath)
	if err != nil {
		return nil, fmt.Errorf("could not read certificate from path '%s': %w", capemPath, err)
	}
	log.Debugf("Certificate read from path: %s", capemPath)

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

	// Deprecated

	// Maybe complete outdated
	// // Collecting integrity information from external service requires nonce
	// // to avoid replay attacks
	// nonce := make([]byte, 8)
	// _, err = rand.Read(nonce)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate nonce: %w", err)
	// }

	// // Connection to CMC
	// ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	// conn, err := grpc.NewClient(d.CmcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Errorf("failed to connect: %w", err)
	// 	cancel()
	// 	return nil, nil
	// }

	// client := grpcapi.NewCMCServiceClient(conn)

	// request := grpcapi.AttestationRequest{
	// 	Nonce: nonce,
	// }

	// // Collect attestation report from CMC
	// response, err := client.Attest(ctx, &request)
	// if err != nil {
	// 	return nil, fmt.Errorf("gRPC Attest call failed: %w", err)
	// }
	// if response.GetStatus() != grpcapi.Status_OK {
	// 	return nil, fmt.Errorf("gRPC Attest call returned status %w", response.GetStatus())
	// }

	// // Verify attestation report
	// result, err := verifyAttestationReport(response.AttestationReport, nonce, capem)
	// if err != nil {
	// 	err = fmt.Errorf("verification failed: %v", err)
	// 	log.Error(err)
	// }

	// r, err := handleReport(result)
	// if err != nil {
	// 	return nil, fmt.Errorf("could not handle attestation report: %w", err)
	// }

	return list, nil
}

// Deprecated
// func verifyAttestationReport(ar, nonce, capem []byte) (ar.VerificationResult, error) {

// 	result := verify.Verify(ar, nonce, capem, nil, verify.PolicyEngineSelect_None, "")
// 	if !result.Success {
// 		return result, fmt.Errorf("verification of attestation report failed")
// 	}

// 	return result, nil
// }

func handleReport(result ar.VerificationResult) (ontology.IsResource, error) {
	raw, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal integrity result: %w", err)
	}

	// TODO(anatheka): Add Resource to Ontology
	resource := &ontology.VirtualMachine{
		Id:   result.Prover,
		Name: result.Prover,
		// CreationTime: , // TODO: TBD
		// GeoLocation: ,// TODO: TBD
		Raw: string(raw),
		RemoteAttestation: &ontology.RemoteAttestation{
			Enabled:      true,
			Status:       result.Success,
			CreationTime: timestamp(result.Created),
		},
	}

	return resource, nil
}

func timestamp(t string) *timestamppb.Timestamp {
	time, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Errorf("could not convert time string to timestamppb: w", err)
		return &timestamppb.Timestamp{}
	}

	return timestamppb.New(time)
}
