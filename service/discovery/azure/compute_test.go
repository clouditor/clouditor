package azure_test

import (
	"net/http"
	"testing"

	"clouditor.io/clouditor/service/discovery/azure"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
)

type mockComputeSender struct {
	mockSender
}

func (m mockComputeSender) Do(req *http.Request) (res *http.Response, err error) {
	var handled bool

	if res, handled, err = m.doSubscriptions(req); handled {
		return
	}

	if req.URL.Path == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Compute/virtualMachines" {
		res, err = createResponse(map[string]interface{}{
			"value": &[]map[string]interface{}{
				{
					"id":         "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Compute/virtualMachines/vm1",
					"name":       "vm1",
					"location":   "eastus",
					"properties": map[string]interface{}{},
				},
			},
		}, 200)
	} else {
		res, err = createResponse(map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)
	}

	return
}

func TestListCompute(t *testing.T) {
	d := azure.NewAzureComputeDiscovery(
		azure.WithSender(&mockComputeSender{}),
		azure.WithAuthorizer(&mockAuthorizer{}),
	)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))

	storage, ok := list[0].(*voc.VirtualMachineResource)

	assert.True(t, ok)
	assert.Equal(t, "vm1", storage.Name)
}
