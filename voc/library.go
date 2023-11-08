// Copyright 2023 Fraunhofer AISEC
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

package voc

var LibraryType = []string{"Library", "Resource"}

// Library is an entity in our Cloud ontology. It encapsulates the (source) code of a library, similar to an application.
// TODO(oxisto): Add this to the ontology and auto-generate it
type Library struct {
	*Resource
	Functionalities     []*Functionality `json:"functionalities"`
	ProgrammingLanguage string           `json:"programmingLanguage"`
	TranslationUnits    []ResourceID     `json:"translationUnits"`
	Dependencies        []ResourceID     `json:"dependencies"`
	GroupID             string           `json:"groupId"`
	ArtifactID          string           `json:"artifactId"`
	Version             string           `json:"version"`
	DependencyType      string           `json:"dependencyType"` // DependencyType denotes which type of dependency it is, e.g., maven or npm
	URL                 string           `json:"url"`
}
