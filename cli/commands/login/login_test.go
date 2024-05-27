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

package login

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"clouditor.io/clouditor/v2/cli"
	"clouditor.io/clouditor/v2/internal/testutil"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	oauth2 "github.com/oxisto/oauth2go"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestLogin(t *testing.T) {
	var (
		err      error
		dir      string
		verifier string
		authSrv  *oauth2.AuthorizationServer
		port     uint16
	)

	authSrv, port, err = testutil.StartAuthenticationServer()
	assert.NoError(t, err)

	dir, err = os.MkdirTemp(os.TempDir(), ".clouditor")
	assert.NoError(t, err)
	assert.NotEmpty(t, dir)

	viper.Set(OAuth2AuthURLFlag, fmt.Sprintf("http://localhost:%d/v1/auth/authorize", port))
	viper.Set(OAuth2TokenURLFlag, fmt.Sprintf("http://localhost:%d/v1/auth/token", port))
	viper.Set(cli.SessionFolderFlag, dir)

	verifier = "012345678901234567890123456789"
	VerifierGenerator = func() string {
		return verifier
	}

	// Issue a code that we can use in the callback
	code := authSrv.IssueCode(oauth2.GenerateCodeChallenge(verifier))

	cmd := NewLoginCommand()

	// Because this potentially blocks, we need to wrap all of this in a timeout
	// TODO(oxisto): This can be removed once https://github.com/golang/go/issues/48157 is fixed
	timeout := time.After(5 * time.Second)
	done := make(chan bool)

	go func() {
		// Simulate a callback
		go func() {
			// Wait for the callback server to be ready
			<-callbackServerReady
			_, err = http.Get(fmt.Sprintf("%s?code=%s", DefaultCallback, code))
			if err != nil {
				assert.NoError(t, err)
			}
		}()

		err = cmd.RunE(nil, []string{"localhost:9090"})
		assert.NoError(t, err)
		done <- true
	}()

	select {
	case <-timeout:
		assert.Fail(t, "Did not finish in time")
	case <-done:
	}
}
