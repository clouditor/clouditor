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
	"context"
	"testing"

	"clouditor.io/clouditor/internal/util"

	"github.com/stretchr/testify/assert"
)

// TODO(lebogg): Add tests for all KeyVault parts here. Try to use mocks provided by Azure

func Test_azureDiscovery_initKeyVaultsClient(t *testing.T) {
	// Init Key Vault Mock and prepare discovery for calling it (TODO(lebogg): Maybe merge these functions)
	initKeyVaultTests()
	d := setDiscoveryForKeyVault()

	err := d.initKeyVaultsClient()

	// Assert if no error was produced and if client is non-empty
	assert.NoError(t, err)
	assert.NotNil(t, d.clients.keyVaultClient)

	// Assert if the Client behaves as expected: Returning exemplary key vault
	get, err := d.clients.keyVaultClient.Get(context.TODO(), "fake-resource-group", "Fake-KeyVault-Name", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Fake-KeyVault-Name", util.Deref(get.Name))
}
