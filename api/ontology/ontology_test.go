package ontology

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestBla(t *testing.T) {
	c := &CloudResource{
		Type: &CloudResource_Compute{
			&Compute{
				Type: &Compute_VirtualMachine{
					&VirtualMachine{
						Name: "my VM",
						BootLogging: &BootLogging{
							Enabled: true,
						},
						ResourceLogging: &ResourceLogging{
							ToId: "my-log-server",
						},
						NetworkInterfaceIds: []string{"test"},
					},
				},
			},
		},
	}
	b, _ := protojson.Marshal(c)
	fmt.Println(string(b))

	n := &CloudResource{
		Type: &CloudResource_Networking{
			&Networking{
				Type: &Networking_NetworkInterface{
					&NetworkInterface{
						Name: "my-id",
						AccessRestriction: &AccessRestriction{
							Type: &AccessRestriction_Firewall{
								&Firewall{
									Type: &Firewall_L3Firewall{
										&L3Firewall{
											Enabled:  true,
											Features: "awesome",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	b, _ = protojson.Marshal(n)
	fmt.Println(string(b))
}
