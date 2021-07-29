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

package voc_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"clouditor.io/clouditor/policies"
	"clouditor.io/clouditor/voc"
)

func TestGraph(t *testing.T) {
	// lets build a graph consisting of a network interface and a VM
	ni := voc.NetworkInterface{
		Networking: &voc.Networking{
			CloudResource: &voc.CloudResource{
				ID: "my-network-interface",
			},
		},
	}

	vm := voc.VirtualMachine{
		Compute: &voc.Compute{
			CloudResource: &voc.CloudResource{
				ID: "my-vm",
			},
		},
	}

	// connect both
	// ni.AttachedTo = vm.ID
	vm.NetworkInterface = append(vm.NetworkInterface, ni.ID)

	out, err := json.Marshal(vm)

	fmt.Printf("%+v\n", err)
	fmt.Printf("%s\n", string(out))

	// make sure, that we are in the clouditor root folder to find the policies
	err = os.Chdir("../")
	if err != nil {
		panic(err)
	}

	tls := map[string]interface{}{
		"enabled": true,
	}

	m := map[string]interface{}{
		"httpEndpoint": map[string]interface{}{
			"transportEncryption": &tls,
		},
	}

	tls["cycle"] = &m

	data, err := policies.RunMap("policies/metric1.rego", m)

	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", data)
}
