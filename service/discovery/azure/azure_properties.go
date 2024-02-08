//go:build exclude

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

	"clouditor.io/clouditor/internal/constants"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
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
func tlsVersion(version string) string {
	// Check TLS version
	switch version {
	case constants.TLS1_0, constants.TLS1_1, constants.TLS1_2:
		return version
	case "1.0", "1_0":
		return constants.TLS1_0
	case "1.1", "1_1":
		return constants.TLS1_1
	case "1.2", "1_2":
		return constants.TLS1_2
	default:
		log.Warningf("'%s' is no implemented TLS version.", version)
		return ""
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

func resourceGroupID(ID *string) voc.ResourceID {
	// split according to "/"
	s := strings.Split(util.Deref(ID), "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 5 {
		return ""
	}

	id := strings.Join(s[:5], "/")

	return voc.ResourceID(id)
}

// retentionDuration returns the retention string as time.Duration
func retentionDuration(retention string) time.Duration {
	if retention == "" {
		return time.Duration(0)
	}

	// Delete first and last character
	r := retention[1 : len(retention)-1]

	// string to int
	d, err := strconv.Atoi(r)
	if err != nil {
		log.Errorf("could not convert string to int")
		return time.Duration(0)
	}

	// Create duration in hours
	duration := time.Duration(time.Duration(d) * time.Hour * 24)

	return duration
}

// labels converts the resource tags to the vocabulary label
func labels(tags map[string]*string) map[string]string {
	l := make(map[string]string)

	for tag, i := range tags {
		l[tag] = util.Deref(i)
	}

	return l
}
