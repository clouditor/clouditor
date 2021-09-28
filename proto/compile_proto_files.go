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

package proto

//go:generate protoc -I ./ -I ../third_party assessment.proto evidence.proto --go_out=../ --go-grpc_out=../ --go_opt=Mevidence.proto=clouditor.io/clouditor/api/evidence --go-grpc_opt=Mevidence.proto=clouditor.io/clouditor/api/evidence --openapi_out=../openapi/assessment
//go:generate protoc -I ./ -I ../third_party auth.proto --go_out=../ --go-grpc_out=../
//go:generate protoc -I ./ -I ../third_party discovery.proto --go_out=../ --go-grpc_out=../ --openapi_out=../openapi/discovery
//go:generate protoc -I ./ -I ../third_party evidence_store.proto evidence.proto --go_out=../ --go-grpc_out=../  --openapi_out=../openapi/evidence
//go:generate protoc -I ./ -I ../third_party orchestrator.proto metric.proto --go_out=../ --go-grpc_out=../ --go_opt=Mmetric.proto=clouditor.io/clouditor/api/assessment --go-grpc_opt=Mmetric.proto=clouditor.io/clouditor/api/assessment --openapi_out=../openapi/orchestrator
