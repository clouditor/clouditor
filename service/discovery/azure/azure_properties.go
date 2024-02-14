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
	"strconv"
	"strings"
	"time"

	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/util"

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

// accountName return the ID's account name
func accountName(id string) string {
	if id == "" {
		return ""
	}

	splitName := strings.Split(id, "/")
	return splitName[8]
}

// tlsVersion returns Clouditor's TLS version constants for the given TLS version
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

func tlsCipherSuites(cs string) []*ontology.CipherSuite {
	// TODO(oxisto): Implement this correctly
	// hack hack
	if cs == "TLS_AES_128_GCM_SHA256" {
		return []*ontology.CipherSuite{
			{
				SessionCipher: "AES-128-GCM",
				MacAlgorithm:  "SHA-256",
			},
		}
	}
	if cs == "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384" {
		return []*ontology.CipherSuite{
			{
				AuthenticationMechanism: "RSA",
				KeyExchangeAlgorithm:    "ECDHE",
				SessionCipher:           "AES-256-GCM",
				MacAlgorithm:            "SHA-384",
			},
		}
	} else {
		return nil
	}
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

	id := strings.Join(s[:5], "/")

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

// labels converts the resource tags to the vocabulary label
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
