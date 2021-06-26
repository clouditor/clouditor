package voc_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"clouditor.io/clouditor/voc"
)

func TestGraph(t *testing.T) {
	// lets build a graph consisting of a network interface and a VM
	ni := voc.NetworkInterface{
		NetworkResource: voc.NetworkResource{
			Resource: voc.Resource{
				ID: "my-network-interface",
			},
		},
	}

	vm := voc.VirtualMachineResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID: "my-vm",
			},
		},
	}

	// connect both
	ni.AttachedTo = vm.ID
	vm.NetworkInterfaces = append(vm.NetworkInterfaces, ni.ID)

	out, err := json.Marshal(vm)

	fmt.Printf("%+v\n", err)
	fmt.Printf("%s\n", string(out))
}
