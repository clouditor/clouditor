package discovery

import (
	"encoding/json"
	"fmt"
	"testing"

	"clouditor.io/clouditor/voc"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test(t *testing.T) {
	infc := Resource{
		Id: "0123123",
		Type: &Resource_CloudResource{
			&CloudResource{
				GeoLocation: &GeoLocation{
					Region: "useast",
				},
				Type: &CloudResource_Networking{
					&Networking{
						Type: &Networking_NetworkInterface{
							&NetworkInterface{
								AccessRestriction: &AccessRestriction{
									Type: &AccessRestriction_Firewall{
										&Firewall{
											Type: &Firewall_L3Firewall{
												&L3Firewall{
													Enabled:         true,
													RestrictedPorts: []string{"80"},
												}},
										}},
								},
							}},
					}},
			}},
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
