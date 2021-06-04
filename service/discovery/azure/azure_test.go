package azure_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "azure-tests")
}

type mockSender struct {
}

func (m mockSender) doSubscriptions(req *http.Request) (res *http.Response, handled bool, err error) {
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
		handled = true
	}

	return
}

type mockAuthorizer struct{}

func (a mockAuthorizer) WithAuthorization() autorest.PrepareDecorator {
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
