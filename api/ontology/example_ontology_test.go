package ontology

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
)

func ExampleMarshal() {
	var vm = &VirtualMachine{
		Id:   "my-id",
		Name: "My VM",
		BootLogging: &BootLogging{
			Enabled: true,
		},
	}

	b, err := protojson.Marshal(vm)
	if err != nil {
		panic(err)
	}

	// We need to use regular JSON marshalling on the output to get a consistent result for tests. See
	// https://github.com/golang/protobuf/issues/1373#issuecomment-946205483
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
	// Output: {"bootLogging":{"enabled":true},"id":"my-id","name":"My VM"}
}
