package discovery

import (
	"encoding/json"
	"fmt"
	"testing"

	"clouditor.io/clouditor/voc"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test(t *testing.T) {
	infc := NetworkInterface{
		Networking: &Networking{
			Resource: &ResourceProperties{
				GeoLocation: &GeoLocation{
					Region: "KOBOLD",
				},
			},
		},
		AccessRestriction: &AccessRestriction{Feature: &AccessRestriction_Firewall{
			Firewall: &Firewall{Class: &Firewall_L3Firewall{
				&L3Firewall{
					Enabled: true, Inbound: false, RestrictedPorts: []string{"80"},
				}},
			}}},
	}

	infc2 := voc.NetworkInterface{
		Networking: &voc.Networking{
			Resource: &voc.Resource{
				GeoLocation: voc.GeoLocation{
					Region: "KOBOLD",
				},
			},
		},
		AccessRestriction: voc.L3Firewall{
			Inbound:         true,
			RestrictedPorts: "80",
		},
	}

	b, _ := protojson.Marshal(&infc)
	fmt.Println(string(b))

	b2, _ := json.Marshal(&infc2)
	fmt.Println(string(b2))
}
