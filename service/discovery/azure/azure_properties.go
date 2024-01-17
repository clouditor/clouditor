package azure

import (
	"strconv"
	"strings"
	"time"

	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
)

// getName returns the name of a given Azure ID
func getName(id string) string {
	if id == "" {
		return ""
	}
	return strings.Split(id, "/")[8]
}

// accountName return the ID's account name
func accountName(id string) string {
	if id == "" {
		return ""
	}

	splitName := strings.Split(id, "/")
	return splitName[8]
}

// generalizeURL generalizes the URL, because the URL depends on the storage type
func generalizeURL(url string) string {
	if url == "" {
		return ""
	}

	urlSplit := strings.Split(url, ".")
	urlSplit[1] = "[file,blob]"
	newURL := strings.Join(urlSplit, ".")

	return newURL
}

// resourceGroupName returns the resource group name of a given Azure ID
func resourceGroupName(id string) string {
	return strings.Split(id, "/")[4]
}

func resourceGroupID(ID *string) voc.ResourceID {
	// split according to "/"
	s := strings.Split(util.Deref(ID), "/")

	// We cannot really return an error here, so we just return an empty string
	if len(s) < 5 {
		return ""
	}

	id := strings.Join(s[:5], "/")

	return voc.ResourceID(id)
}

// retentionDuration returns the retention string as time.Duration
func retentionDuration(retention string) time.Duration {
	if retention == "" {
		return time.Duration(0)
	}

	// Delete first and last character
	r := retention[1 : len(retention)-1]

	// string to int
	d, err := strconv.Atoi(r)
	if err != nil {
		log.Errorf("could not convert string to int")
		return time.Duration(0)
	}

	// Create duration in hours
	duration := time.Duration(time.Duration(d) * time.Hour * 24)

	return duration
}

// labels converts the resource tags to the vocabulary label
func labels(tags map[string]*string) map[string]string {
	l := make(map[string]string)

	for tag, i := range tags {
		l[tag] = util.Deref(i)
	}

	return l
}
