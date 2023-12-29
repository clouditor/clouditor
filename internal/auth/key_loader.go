// Copyright 2016-2020 Fraunhofer AISEC
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

package auth

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	// DefaultApiKeySaveOnCreate specifies whether a created API key will be saved. This is useful to turn of in unit tests, where
	// we only want a temporary key.
	DefaultApiKeySaveOnCreate = true

	// DefaultApiKeyPassword is the default password to protect the API key
	DefaultApiKeyPassword = "changeme"

	// DefaultApiKeyPath is the default path for the API private key
	DefaultApiKeyPath = DefaultConfigDirectory + "/api.key"

	// DefaultConfigDirectory is the default path for the clouditor configuration, such as keys
	DefaultConfigDirectory = "~/.clouditor"
)

// ExpandPath expands a path that possible contains a tilde (~) character into the home directory
// of the user
func ExpandPath(path string) (out string, err error) {
	var (
		u *user.User
	)

	// Fetch the current user home directory
	u, err = user.Current()
	if err != nil {
		return path, fmt.Errorf("could not find retrieve current user: %w", err)
	}

	if path == "~" {
		return u.HomeDir, nil
	} else if strings.HasPrefix(path, "~") {
		// We only allow ~ at the beginning of the path
		return filepath.Join(u.HomeDir, path[2:]), nil
	}

	return path, nil
}
