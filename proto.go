/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package clouditor

// Run owl2proto to generate the ontology proto file. The tools needs as arguments the following:
// - root-resource-name from the owl file
// - owl file in the owx format
// - header-file
// - output-path for the proto file (optional, default is "api/ontology.proto")
//go:generate owl2proto generate-proto --root-resource-name=https://ontology.emerald-he.eu/classes/Resource internal/ontology/ontology.owx --header-file=internal/ontology/clouditor_header.proto --output-path=api/ontology/ontology.proto --full-semantic-mode=false --deterministic-field-numbers=true
//go:generate buf format -w
//go:generate buf generate --exclude-path="internal/ontology/clouditor_header.proto"
//go:generate buf generate --exclude-path="internal/ontology/clouditor_header.proto" --template buf.gotag.gen.yaml
//go:generate buf generate --template buf.openapi.gen.yaml --path api/assessment -o openapi/assessment
//go:generate buf generate --template buf.openapi.gen.yaml --path api/evaluation -o openapi/evaluation
//go:generate buf generate --template buf.openapi.gen.yaml --path api/discovery -o openapi/discovery
//go:generate buf generate --template buf.openapi.gen.yaml --path api/evidence -o openapi/evidence
//go:generate buf generate --template buf.openapi.gen.yaml --path api/orchestrator -o openapi/orchestrator
//go:generate buf generate --template buf.openapi.gen.yaml --path api/ontology -o openapi/ontology
