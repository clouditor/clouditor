// Copyright 2022 Fraunhofer AISEC
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

package defaults

const (
	DefaultTargetCloudServiceID          = "00000000-0000-0000-0000-000000000000"
	DefaultTargetCloudServiceName        = "default"
	DefaultTargetCloudServiceDescription = "The default target cloud service"
	DefaultCatalogID                     = "EUCS"
	DefaultEUCSLowerLevelControlID136    = "OPS-13.6"
	DefaultEUCSLowerLevelControlID137    = "OPS-13.7"
	DefaultEUCSUpperLevelControlID13     = "OPS-13"
	DefaultEUCSCategoryName              = "Operational Security"
	DefaultOrchestratorAddress           = "localhost:9090"
)

var (
	AssuranceLevelHigh   = "high"
	AssuranceLevelMedium = "medium"
	AssuranceLevelBasic  = "basic"
)
