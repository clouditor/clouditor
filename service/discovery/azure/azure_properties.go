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

package azure

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// getName returns the name of a given Azure ID
func getName(id string) string {
	if id == "" {
		return ""
	}

	return strings.Split(id, "/")[8]
}

// resourceID makes sure that the Azure ID we get is lowercase, because Azure sometimes has weird notions that things
// are uppercase. Their documentation says that comparison of IDs is case-insensitive, so we lowercase everything.
func resourceID(id *string) string {
	if id == nil {
		return ""
	}

	return strings.ToLower(*id)
}

// resourceIDPointer makes sure that the Azure ID we get is lowercase, because Azure sometimes has weird notions that things are uppercase. Their documentation says that comparison of IDs is case-insensitive, so we lowercase everything.
func resourceIDPointer(id *string) *string {
	return util.Ref(resourceID(id))
}

// accountName return the ID's account name
func accountName(id string) string {
	if id == "" {
		return ""
	}

	splitName := strings.Split(id, "/")
	return splitName[8]
}

// tlsVersion returns a float value for the given TLS version string
func tlsVersion(version *string) float32 {
	if version == nil {
		return 0
	}

	// Check TLS version
	switch *version {
	case "1.0", "1_0", string(armstorage.MinimumTLSVersionTLS10):
		return 1.0
	case "1.1", "1_1", string(armstorage.MinimumTLSVersionTLS11):
		return 1.1
	case "1.2", "1_2", string(armstorage.MinimumTLSVersionTLS12):
		return 1.2
	case "1.3", "1_3":
		return 1.3
	default:
		log.Warningf("'%s' is not an implemented TLS version.", *version)
		return 0
	}
}

// tlsCipherSuites parses TLS cipher suites. Examples are TLS_AES_128_GCM_SHA256 or
// TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384.
func tlsCipherSuites(cs string) []*ontology.CipherSuite {
	var (
		parts  []string
		i      int
		cipher ontology.CipherSuite
	)

	parts = strings.Split(cs, "_")

	if parts[i] != "TLS" {
		return nil
	}

	// Next is either a key exchange or directly the session cipher
	i++
	if parts[i] == "ECDHE" {
		cipher.KeyExchangeAlgorithm = parts[i]
	} else {
		i--
		goto cipher
	}

	i++
	if slices.Contains([]string{"RSA", "ECDSA"}, parts[i]) {
		cipher.AuthenticationMechanism = parts[i]
	} else {
		goto invalid
	}

	i++
	if parts[i] != "WITH" {
		goto invalid
	}

cipher:
	i++
	if parts[i] == "AES" {
		cipher.SessionCipher = parts[i]
	} else {
		goto invalid
	}

	i++
	if slices.Contains([]string{"128", "256"}, parts[i]) {
		cipher.SessionCipher += "-" + parts[i]
	} else {
		goto invalid
	}

	i++
	if slices.Contains([]string{"CBC", "GCM"}, parts[i]) {
		cipher.SessionCipher += "-" + parts[i]
	} else {
		goto invalid
	}

	i++
	if parts[i] == "SHA256" {
		cipher.MacAlgorithm = "SHA-256"
	} else if parts[i] == "SHA384" {
		cipher.MacAlgorithm = "SHA-384"
	} else {
		goto invalid
	}

	return []*ontology.CipherSuite{&cipher}

invalid:
	return nil
}

// generalizeURL generalizes the URL, because the URL depends on the storage type
func generalizeURL(url string) string {
	if url == "" {
		return ""
	}

	urlSplit := strings.Split(url, ".")
	urlSplit[1] = "[file,blob]"
	newURL := strings.Join(urlSplit, ".")

	return newURL
}

// resourceGroupName returns the resource group name of a given Azure ID
func resourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}

// resourceGroupID builds a resource group ID out of the resource ID. It will also be lowercase (see resourceID function
// for reasoning).
func resourceGroupID(ID *string) *string {
	if ID == nil {
		return nil
	}

	// split according to "/"
	s := strings.Split(util.Deref(ID), "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 5 {
		return nil
	}

	id := strings.ToLower(strings.Join(s[:5], "/"))

	return &id
}

// retentionDuration returns the retention string as time.Duration
func retentionDuration(retention string) *durationpb.Duration {
	if retention == "" {
		return durationpb.New(time.Duration(0))
	}

	// Delete first and last character
	r := retention[1 : len(retention)-1]

	// string to int
	d, err := strconv.Atoi(r)
	if err != nil {
		log.Errorf("could not convert string to int")
		return durationpb.New(time.Duration(0))
	}

	// Create duration in hours
	duration := time.Duration(time.Duration(d) * time.Hour * 24)

	return durationpb.New(duration)
}

// labels converts the resource tags to the ontology label
func labels(tags map[string]*string) map[string]string {
	l := make(map[string]string)

	for tag, i := range tags {
		l[tag] = util.Deref(i)
	}

	return l
}

func creationTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func location(region *string) *ontology.GeoLocation {
	if region == nil {
		return nil
	}

	return &ontology.GeoLocation{
		Region: util.Deref(region),
	}
}

// publicNetworkAccessStatus returns the public access status of the resource
func publicNetworkAccessStatus(status *string) bool {
	// Check if a mandatory field is empty
	if status == nil {
		return false
	}

	// Check if resource is public available
	if util.Deref(status) == "Enabled" {
		return true
	}

	return false
}

// discoverDiagnosticSettings discovers the diagnostic setting for the given resource URI and returns the information of the needed information of the log properties as ontology.ActivityLogging object and the Azure response.
func (d *azureDiscovery) discoverDiagnosticSettings(resourceURI string) (*ontology.ActivityLogging, string, error) {
	var (
		al           *ontology.ActivityLogging
		workspaceIDs []string
		raw          string
	)

	if err := d.initDiagnosticsSettingsClient(); err != nil {
		return nil, "", err
	}

	// List all diagnostic settings for the storage account
	listPager := d.clients.diagnosticSettingsClient.NewListPager(resourceURI, &armmonitor.DiagnosticSettingsClientListOptions{})
	for listPager.More() {
		pageResponse, err := listPager.NextPage(context.TODO())
		if err != nil {
			err = fmt.Errorf("%s: %v", ErrGettingNextPage, err)
			return nil, "", err
		}

		for _, value := range pageResponse.Value {
			// Check if data is sent to a log analytics workspace
			if value.Properties.WorkspaceID == nil {
				log.Debugf("diagnostic setting '%s' does not send data to a Log Analytics Workspace", util.Deref(value.Name))
				continue
			}

			// Add Log Analytics WorkspaceIDs to slice
			workspaceIDs = append(workspaceIDs, util.Deref(value.Properties.WorkspaceID))
		}

		raw = discovery.Raw(pageResponse)
	}

	if len(workspaceIDs) > 0 {
		al = &ontology.ActivityLogging{
			Enabled:           true,
			LoggingServiceIds: workspaceIDs, // TODO(all): Each diagnostic setting has also a retention period, maybe we should add that information as well
		}
	}

	return al, raw, nil
}
