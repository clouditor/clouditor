package ontology

import (
	"encoding/json"
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

	// We need to use regular JSON marshalling on the output to get a consistent result for tests
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	b, err = json.Marshal(&m)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	// Output: {"cloudResource":{"compute":{"virtualMachine":{"bootLogging":{"enabled":true},"id":"my-id","name":"My VM"}}}}
}
