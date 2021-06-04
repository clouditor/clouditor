package azure_test

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

type mocky struct{}

func (m mocky) Do(req *http.Request) (res *http.Response, err error) {
	if req.URL.Path == "/subscriptions" {
		res, err = createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":             "/subscriptions/00000000-0000-0000-0000-000000000000",
					"subscriptionId": "00000000-0000-0000-0000-000000000000",
					"name":           "sub1",
				},
			},
		}, 200)

		return
	}

	res, err = createResponse(map[string]interface{}{
		"value": &[]map[string]interface{}{
			{
				"id":       "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account1",
				"name":     "account1",
				"location": "eastus",
				"properties": map[string]interface{}{
					"creationTime": "2017-05-24T13:28:53.4540398Z",
					"primaryEndpoints": map[string]interface{}{
						"blob": "https://account1.blob.core.windows.net/",
					},
					"encryption": map[string]interface{}{
						"services": map[string]interface{}{
							"file": map[string]interface{}{
								"keyType":         "Account",
								"enabled":         true,
								"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
							},
							"blob": map[string]interface{}{
								"keyType":         "Account",
								"enabled":         true,
								"lastEnabledTime": "2019-12-11T20:49:31.7036140Z",
							},
						},
						"keySource": "Microsoft.Storage",
					},
				},
			},
		},
	}, 200)

	return
}

func TestListStorage(t *testing.T) {
	d := azure.NewAzureStorageDiscovery(azure.WithSender(&mocky{}), azure.WithAuthorizer(&mockAuthorizer{}))

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	storage, ok := list[0].(*voc.ObjectStorageResource)

	assert.True(t, ok)
	assert.Equal(t, "account1", storage.Name)
}
