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
