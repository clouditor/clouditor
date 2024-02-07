package ontology

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
)

func ExampleMarshal() {
	var c = &Resource{
		Type: &Resource_CloudResource{&CloudResource{
			Type: &CloudResource_Compute{Compute: &Compute{
				Type: &Compute_VirtualMachine{&VirtualMachine{
					Id:   "my-id",
					Name: "My VM",
					BootLogging: &BootLogging{
						Enabled: true,
					},
				}},
			}},
		}},
	}

	b, err := protojson.Marshal(c)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	// Output: {"cloudResource":{"compute":{"virtualMachine":{"id":"my-id", "name":"My VM", "bootLogging":{"enabled":true}}}}}
}
