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
	"time"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewKubernetesStorageDiscovery(t *testing.T) {
	type args struct {
		intf                  kubernetes.Interface
		CertificationTargetID string
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "empty input",
			want: &k8sStorageDiscovery{
				k8sDiscovery: k8sDiscovery{},
			},
		},
		{
			name: "Happy path",
			args: args{
				intf:                  &fake.Clientset{},
				CertificationTargetID: testdata.MockCertificationTargetID1,
			},
			want: &k8sStorageDiscovery{
				k8sDiscovery: k8sDiscovery{
					intf: &fake.Clientset{},
					ctID: testdata.MockCertificationTargetID1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKubernetesStorageDiscovery(tt.args.intf, tt.args.CertificationTargetID)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
			assert.Equal(t, "Kubernetes Storage", got.Name())
		})
	}
}

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

	d := NewKubernetesStorageDiscovery(client, testdata.MockCertificationTargetID1)

	list, err := d.List()
	assert.NoError(t, err)
	assert.NotNil(t, list)

	// Check persistentVolume
	volume, ok := list[0].(*ontology.BlockStorage)

	// Create expected ontology.BlockStorage
	expectedVolume := &ontology.BlockStorage{
		Id:               volumeUID,
		Name:             volumeName,
		CreationTime:     volume.CreationTime,
		Labels:           volumeLabel,
		AtRestEncryption: &ontology.AtRestEncryption{},
	}

	assert.True(t, ok)
	assert.Equal(t, expectedVolume, volume, protocmp.IgnoreFields(&ontology.BlockStorage{}, "raw"))
}

func Test_k8sStorageDiscovery_handlePV(t *testing.T) {
	type fields struct {
		k8sDiscovery k8sDiscovery
	}
	type args struct {
		pv *corev1.PersistentVolume
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   ontology.IsResource
	}{
		{
			name:   "file-based",
			fields: fields{},
			args: args{
				pv: &corev1.PersistentVolume{
					ObjectMeta: metav1.ObjectMeta{
						UID:               "my-id",
						Name:              "test",
						CreationTimestamp: metav1.NewTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					},
					Spec: corev1.PersistentVolumeSpec{
						PersistentVolumeSource: corev1.PersistentVolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/tmp",
							},
						},
					},
				},
			},
			want: &ontology.FileStorage{
				Id:               "my-id",
				Name:             "test",
				AtRestEncryption: &ontology.AtRestEncryption{},
				CreationTime:     timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				Raw:              `{"*v1.PersistentVolume":[{"metadata":{"name":"test","uid":"my-id","creationTimestamp":"2024-01-01T00:00:00Z"},"spec":{"hostPath":{"path":"/tmp"}},"status":{}}]}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &k8sStorageDiscovery{
				k8sDiscovery: tt.fields.k8sDiscovery,
			}
			got := d.handlePV(tt.args.pv)
			assert.Equal(t, tt.want, got)
		})
	}
}
