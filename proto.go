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

//go:generate protoc -I . -I third_party internal/testutil/prototest/prototest.proto --go_out=paths=source_relative:.
//go:generate protoc -I . -I third_party api/page_token.proto --go_out=paths=source_relative:.
//go:generate protoc -I . -I third_party api/assessment/metric.proto --go_out=paths=source_relative:.
//go:generate protoc -I . -I third_party api/evidence/evidence.proto --go_out=paths=source_relative:.
//go:generate protoc -I . -I third_party api/assessment/assessment.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --openapi_out=openapi/assessment --openapi_opt enum_type=string --grpc-gateway_out=paths=source_relative:. --grpc-gateway_opt logtostderr=true
//go:generate protoc -I . -I third_party api/auth/auth.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=paths=source_relative:. --grpc-gateway_opt logtostderr=true
//go:generate protoc -I . -I third_party api/discovery/discovery.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --openapi_out=openapi/discovery --openapi_opt enum_type=string --grpc-gateway_out=paths=source_relative:. --grpc-gateway_opt logtostderr=true
//go:generate protoc -I . -I third_party api/evidence/evidence_store.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --openapi_out=openapi/evidence --openapi_opt enum_type=string --grpc-gateway_out=paths=source_relative:. --grpc-gateway_opt logtostderr=true
//go:generate protoc -I . -I third_party api/orchestrator/orchestrator.proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --openapi_out=openapi/orchestrator --openapi_opt enum_type=string --grpc-gateway_out=paths=source_relative:. --grpc-gateway_opt logtostderr=true
//go:generate protoc -I . -I third_party --gotag_out=paths=source_relative:. --gotag_opt=Mapi/orchestrator/orchestrator.proto=clouditor.io/api/orchestrator api/orchestrator/orchestrator.proto --gotag_opt=Mapi/assessment/metric.proto=clouditor.io/api/assessment api/assessment/metric.proto --gotag_opt=Mapi/assessment/assessment.proto=clouditor.io/api/assessment api/assessment/assessment.proto
