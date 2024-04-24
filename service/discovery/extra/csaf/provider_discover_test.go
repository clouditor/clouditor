package csaf

import (
	"net/http"
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
)

func Test_csafDiscovery_discoverProviders(t *testing.T) {
	type fields struct {
		domain string
		csID   string
		client *http.Client
	}
	type args struct {
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantProviders []ontology.IsResource
		wantErr       bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &csafDiscovery{
				domain: tt.fields.domain,
				csID:   tt.fields.csID,
				client: tt.fields.client,
			}
			gotProviders, err := d.discoverProviders()
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotProviders, tt.wantProviders) {
				t.Errorf("discoverProviders() gotProviders = %v, want %v", gotProviders, tt.wantProviders)
			}
		})
	}
}
