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

package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httputil"
	"testing"
)

func init() {
	log = logrus.WithField("component", "azure-tests")
}

type mockSender struct {
}

func (mockSender) Do(req *http.Request) (res *http.Response, err error) {
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
	} else {
		res, err = createResponse(map[string]interface{}{}, 404)
		log.Errorf("Not handling mock for %s yet", req.URL.Path)
	}

	return
}

type mockAuthorizer struct{}

func (mockAuthorizer) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return p
	}
}

func createResponse(object map[string]interface{}, statusCode int) (res *http.Response, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)

	if err = enc.Encode(object); err != nil {
		return nil, fmt.Errorf("could not encode JSON object: %w", err)
	}

	body := io.NopCloser(buf)

	return &http.Response{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func LogRequest() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)

			if err != nil {
				log.Println(err)
			}

			dump, _ := httputil.DumpRequestOut(r, true)
			log.Println(string(dump))

			return r, err
		})
	}
}

func LogResponse() autorest.RespondDecorator {
	return func(p autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(r *http.Response) error {
			err := p.Respond(r)

			if err != nil {
				log.Println(err)
			}

			dump, _ := httputil.DumpResponse(r, true)
			log.Println(string(dump))

			return err
		})
	}
}

func TestGetResourceGroupName(t *testing.T) {
	accountId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/res1/providers/Microsoft.Storage/storageAccounts/account3"
	result := getResourceGroupName(accountId)

	assert.Equal(t, "res1", result)
}

func TestApply(t *testing.T) {

	// Test senderOption
	so := senderOption{
		sender: mockStorageSender{},
	}

	client := autorest.Client{}
	so.apply(&client)
	assert.Equal(t, so.sender, client.Sender)

	// Test authorizerOption
	ao := authorizerOption{
		authorizer: mockAuthorizer{},
	}

	ao.apply(&client)
	assert.Equal(t, ao.authorizer, client.Authorizer)

	// Test azureDiscovery
	ad := azureDiscovery{
		authOption: &authorizerOption{
			authorizer: mockAuthorizer{},
		},
	}

	ad.apply(&client)
	assert.Equal(t, ad.authOption.authorizer, client.Authorizer)
}

func TestWithSender(t *testing.T) {
	expected := &senderOption{
		sender: mockStorageSender{},
	}

	resp := WithSender(mockStorageSender{})

	assert.Equal(t, expected, resp)
}

func TestWithAuthorizer(t *testing.T) {
	expected := &authorizerOption{
		authorizer: mockAuthorizer{},
	}

	resp := WithAuthorizer(mockAuthorizer{})

	assert.Equal(t, expected, resp)
}
