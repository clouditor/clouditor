// Copyright 2021 Fraunhofer AISEC
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

package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"os"

	"clouditor.io/clouditor/api/orchestrator"
)

// TODO(oxisto): actually persist them
var accounts []*orchestrator.CloudAccount

func init() {
	// initialize accounts if the directories exists for now, we need to
	// replace this in the future
	if existsInHome(".kube") {
		accounts = append(accounts, &orchestrator.CloudAccount{AccountType: "k8s", AutoDiscover: true})
	}

	if existsInHome(".aws") {
		accounts = append(accounts, &orchestrator.CloudAccount{AccountType: "aws", AutoDiscover: true})
	}

	if existsInHome(".Azure") {
		accounts = append(accounts, &orchestrator.CloudAccount{AccountType: "azure", AutoDiscover: true})
	}
}

func existsInHome(path string) bool {
	home, err := os.UserHomeDir()
	if err == nil {
		_, err := os.Stat(fmt.Sprintf("%s/%s", home, path))
		if err == nil {
			return true
		}
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}

	return false
}

func (*Service) CreateAccount(_ context.Context, req *orchestrator.CreateAccountRequest) (response *orchestrator.CloudAccount, err error) {
	response = req.Account

	if req.Account.AutoDiscover {
		// TODO: do some actual auto discovering
		log.Debugf("Trying to auto-discover an %s account", req.Account.AccountType)
	}

	accounts = append(accounts, req.Account)

	return response, nil
}

func (*Service) ListAccounts(_ context.Context, _ *orchestrator.ListAccountsRequests) (response *orchestrator.ListAccountsResponse, err error) {
	response = new(orchestrator.ListAccountsResponse)

	response.Accounts = accounts

	return response, nil
}
