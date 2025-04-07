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

package testdata

const (
	// auth options
	MockOpenstackIdentityEndpoint = "https://identityHost:portNumber/v2.0" //"https://openstack.test:8888/v2.0"
	MockOpenstackUsername         = "username"
	MockOpenstackPassword         = "password"
	MockOpenstackTenantName       = "openstackTenant"

	// domain
	MockDomainID1          = "00000000000000000000000000000000"
	MockDomainName1        = "Domain 1"
	MockDomainDescription1 = "This is a mock domain description (1)"

	// project
	MockProjectID1          = "00000000000000000000000000000001"
	MockProjectName1        = "Project 1"
	MockProjectDescription1 = "This is a mock project description (1)."
	MockProjectParentID1    = MockDomainID1

	// server
	MockServerID1      = "00000000000000000000000000000002"
	MockServerName1    = "Server 1"
	MockServerTenantID = MockDomainID1

	// volume
	MockVolumeID1      = "00000000000000000000000000000003"
	MockVolumeName1    = "Volume 1"
	MockVolumeTenantID = MockDomainID1

	// network
	MockNetworkID1   = "00000000000000000000000000000004"
	MockNetworkName1 = "Network 1"
)
