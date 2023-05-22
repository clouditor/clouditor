// Copyright 2022 Fraunhofer AISEC
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

	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/voc"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_k8sStorageDiscovery_List(t *testing.T) {

	var (
		volumeName              = "my-volume"
		volumeUID               = "00000000-0000-0000-0000-000000000000"
		volumeCreationTime      = metav1.Now()
		volumeLabel             = map[string]string{"my": "label"}
		persistenVolumeDiskName = "my-disk"
	)

	client := fake.NewSimpleClientset()

	// Create persistent volumes
	v := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:              volumeName,
			UID:               types.UID(volumeUID),
			CreationTimestamp: volumeCreationTime,
			Labels:            volumeLabel,
		},
		Spec: corev1.PersistentVolumeSpec{
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				AzureDisk: &corev1.AzureDiskVolumeSource{
					DiskName: persistenVolumeDiskName,
				},
			},
		},
	}

	_, err := client.CoreV1().PersistentVolumes().Create(context.TODO(), v, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting volume add: %v", err)
	}

	d := NewKubernetesStorageDiscovery(client, testdata.MockCloudServiceID1)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	// Check persistentVolume
	volume, ok := list[0].(*voc.BlockStorage)

	// Create exptected voc.BlockStorage
	expectedVolume := &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(volumeUID),
				ServiceID:    testdata.MockCloudServiceID1,
				Name:         volumeName,
				CreationTime: volume.CreationTime,
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: "",
				},
				Labels: volumeLabel,
			},
			AtRestEncryption: &voc.AtRestEncryption{},
		},
	}

	// Delete raw. We have to delete it, because of the creation time included in the raw field.
	assert.NotNil(t, volume.Raw)
	volume.Raw = ""

	assert.True(t, ok)
	assert.Equal(t, expectedVolume, volume)
}
