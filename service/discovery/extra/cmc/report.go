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
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"clouditor.io/clouditor/v2/api/ontology"

	"github.com/Fraunhofer-AISEC/cmc/attestationreport"
	ci "github.com/Fraunhofer-AISEC/cmc/cmcinterface" // TODO: What do we need here? I think that changed.

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	timeoutSec = 10
)

// discoverReports discovers the attestation reports from the CMC
func (d *cmcDiscovery) discoverReports() ([]ontology.IsResource, error) {
	var (
		list []ontology.IsResource
		// capem = []byte(rawConfig.Certificate)
	)

	// Collecting integrity information from external service requires nonce
	// to avoid replay attacks
	nonce := make([]byte, 8)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Connection to CMC
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	conn, err := grpc.NewClient(d.CmcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("failed to connect: %w", err)
		cancel()
		return nil, nil
	}

	client := ci.NewCMCServiceClient(conn)

	request := ci.AttestationRequest{
		Nonce: nonce,
	}

	// Collect attestation report from CMC
	response, err := client.Attest(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("gRPC Attest call failed: %w", err)
	}
	if response.GetStatus() != ci.Status_OK {
		return nil, fmt.Errorf("gRPC Attest call returned status %w", response.GetStatus())
	}

	// Verify attestation report
	result, err := verifyAttestationReport(response.AttestationReport, nonce, capem) //TODO: What capem do we need here?
	if err != nil {
		err = fmt.Errorf("verification failed: %v", err)
		log.Error(err)
	}

	r, err := handleReport(result)
	if err != nil {
		return nil, fmt.Errorf("could not handle attestation report: %w", err)
	}

	list = append(list, r)

	return list, nil
}

func verifyAttestationReport(ar, nonce, capem []byte) (attestationreport.VerificationResult, error) {

	result := attestationreport.Verify(string(ar), nonce, capem, nil)
	if !result.Success {
		return result, fmt.Errorf("verification of attestation report failed")
	}

	return result, nil
}

func handleReport(result attestationreport.VerificationResult) (ontology.IsResource, error) {
	rawEvidence, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal integrity result: %w", err)
	}

	resource := &ontology.VirtualMachine{
		// Id:   requestId, //TODO: Is there any ID? IP or something else?
		// Name: requestId, //TODO: Is there any name? IP or something else?
		// CreationTime: , // TODO: TBD
		// GeoLocation: ,// TODO: TBD

		Raw: string(rawEvidence),
		// TargetService:  req.ServiceId,
		// TargetResource: result.PlainAttReport.DeviceDescription.Fqdn,
		// ToolId:         ComponentID,
		// GatheredAt:     timestamppb.Now(),
		// Value:          evidenceValue,
		// RawEvidence:    string(rawEvidence),
	}

	return resource, nil
}
