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
	"clouditor.io/clouditor/internal/util"
	"context"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2/fake"
)

var FakeFarmServer = fake.PlansServer{}

// init initializes fake functions, e.g. for fake Plans client
func init() {
	initFakeFarmServer()
	// Web App and Function Server will follow
}

func initFakeFarmServer() {
	// Set up fake Get Function for Farms
	FakeFarmServer.Get = func(ctx context.Context, resourceGroupName string, farmName string,
		options *armappservice.PlansClientGetOptions) (
		resp azfake.Responder[armappservice.PlansClientGetResponse], errResp azfake.ErrorResponder) {
		if farmName == "FarmWithZoneRedundancy" {
			resp.SetResponse(200, armappservice.PlansClientGetResponse{
				Plan: armappservice.Plan{
					Properties: &armappservice.PlanProperties{
						ZoneRedundant: util.Ref(true),
					},
				},
			}, &azfake.SetResponseOptions{})
		}
		if farmName == "FarmWithNoRedundancy" {
			resp.SetResponse(200, armappservice.PlansClientGetResponse{
				Plan: armappservice.Plan{
					Properties: &armappservice.PlanProperties{
						ZoneRedundant: util.Ref(false),
					},
				},
			}, &azfake.SetResponseOptions{})
		}
		return
	}
}
