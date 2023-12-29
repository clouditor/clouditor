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

package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetJSONFilenames returns all json files in the given folder
func GetJSONFilenames(folder string) ([]string, error) {
	var (
		list []string
	)

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for i := range files {
		if strings.HasSuffix(files[i].Name(), ".json") {
			list = append(list, fmt.Sprintf("%s/%s", folder, files[i].Name()))
		}
	}

	return list, nil
}

// userHomeDirFunc points to a function that returns the user home directory. This can be changed for mock tests.
var userHomeDirFunc = os.UserHomeDir

// ExpandPath expands a path that possible contains a tilde (~) character into the home directory
// of the user
func ExpandPath(path string) (out string, err error) {
	var (
		home  string
		found bool
	)

	// Fetch the current user home directory
	home, err = userHomeDirFunc()
	if err != nil {
		return "", fmt.Errorf("could not find retrieve current user: %w", err)
	}

	out, found = strings.CutPrefix(path, "~")
	if found {
		out = filepath.Join(home, out)
	}

	return
}
