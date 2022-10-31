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

package k8s

import (
	"context"
	"testing"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListPods(t *testing.T) {

	var (
		volumeName      = "my-volume"
		diskName        = "my-disk"
		podCreationTime = metav1.Now()
		podName         = "my-pod"
		podID           = "/namespaces/my-namespace/containers/my-pod"
		podNamespace    = "my-namespace"
		podLabel        = map[string]string{"my": "label"}
	)

	client := fake.NewSimpleClientset()

	// Create an Pod with name, creationTimestamp and a AzureDisk volume
	p := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:              podName,
			CreationTimestamp: podCreationTime,
			Labels:            podLabel,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						AzureDisk: &corev1.AzureDiskVolumeSource{
							DiskName: diskName,
						},
					},
				},
			},
		},
	}
	_, err := client.CoreV1().Pods(podNamespace).Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	d := NewKubernetesComputeDiscovery(client)

	list, err := d.List()

	assert.NoError(t, err)
	assert.NotNil(t, list)

	// Check container
	container, ok := list[0].(*voc.Container)

	// Create expected voc.Container
	expectedContainer := &voc.Container{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(podID),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         podName,
				CreationTime: podCreationTime.Unix(),
				Type:         []string{"Container", "Compute", "Resource"},
				Labels:       podLabel,
			},
			NetworkInterfaces: []voc.ResourceID{
				voc.ResourceID(podNamespace),
			},
		},
	}

	assert.True(t, ok)
	assert.Equal(t, expectedContainer, container)

	// Check volume
	volume, ok := list[1].(*voc.BlockStorage)
	// Create expected voc.BlockStorage
	expectedVolume := &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(volumeName),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         volumeName,
				CreationTime: 0,
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: "",
				},
			},
			AtRestEncryption: &voc.AtRestEncryption{},
		},
	}

	assert.True(t, ok)
	assert.Equal(t, expectedVolume, volume)
}
