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

//go:build exclude

package azure

import (
	"strings"

	"clouditor.io/clouditor/voc"
)

// backupsEmptyCheck checks if the backups list is empty and returns voc.Backup with enabled = false.
func backupsEmptyCheck(backups []*voc.Backup) []*voc.Backup {
	if len(backups) == 0 {
		return []*voc.Backup{
			{
				Enabled:         false,
				RetentionPeriod: -1,
				Interval:        -1,
			},
		}
	}

	return backups
}

// backupPolicyName returns the backup policy name of a given Azure ID
func backupPolicyName(id string) string {
	// split according to "/"
	s := strings.Split(id, "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 10 {
		return ""
	}
	return s[10]
}
