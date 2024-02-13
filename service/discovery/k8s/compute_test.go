//go:build exclude

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
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/ontology"
	"clouditor.io/clouditor/internal/testdata"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewKubernetesComputeDiscovery(t *testing.T) {
	type args struct {
		intf           kubernetes.Interface
		cloudServiceID string
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			name: "empty input",
			want: &k8sComputeDiscovery{
				k8sDiscovery: k8sDiscovery{},
			},
		},
		{
			name: "Happy path",
			args: args{
				intf:           &fake.Clientset{},
				cloudServiceID: testdata.MockCloudServiceID1,
			},
			want: &k8sComputeDiscovery{
				k8sDiscovery: k8sDiscovery{
					intf: &fake.Clientset{},
					csID: testdata.MockCloudServiceID1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKubernetesComputeDiscovery(tt.args.intf, tt.args.cloudServiceID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubernetesComputeDiscovery() = %v, want %v", got, tt.want)
			}

			assert.Equal(t, "Kubernetes Compute", got.Name())
		})
	}
}

func Test_k8sComputeDiscovery_List(t *testing.T) {
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

	type fields struct {
		discovery discovery.Discoverer
	}
	tests := []struct {
		name    string
		fields  fields
		want    assert.ValueAssertionFunc
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				NewKubernetesComputeDiscovery(client, testdata.MockCloudServiceID1),
			},
			want: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				cloudResource, ok := i1.([]ontology.IsResource)
				if !assert.True(tt, ok) {
					return false
				}
				container, ok := cloudResource[0].(*voc.Container)
				if !assert.True(tt, ok) {
					return false
				}
				// Create expected voc.Container
				expectedContainer := &voc.Container{
					Compute: &voc.Compute{
						Resource: &voc.Resource{
							ID:        voc.ResourceID(podID),
							ServiceID: testdata.MockCloudServiceID1,
							Name:      podName,
							Type:      []string{"Container", "Compute", "Resource"},
							Labels:    podLabel,
							Raw:       "",
						},
						NetworkInterfaces: []voc.ResourceID{
							voc.ResourceID(podNamespace),
						},
					},
				}

				// Delete creation time
				assert.NotNil(t, container.CreationTime)
				container.CreationTime = 0

				// Delete raw. We have to delete it, because of the creation time included in the raw field.
				assert.NotNil(t, container.Raw)
				container.Raw = ""

				assert.True(t, ok)
				assert.Equal(t, expectedContainer, container)

				// Check volume
				volume, ok := cloudResource[1].(*voc.BlockStorage)
				// Create expected voc.BlockStorage
				expectedVolume := &voc.BlockStorage{
					Storage: &voc.Storage{
						Resource: &voc.Resource{
							ID:           voc.ResourceID(volumeName),
							ServiceID:    testdata.MockCloudServiceID1,
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

				// Delete raw. We have to delete it, because of the creation time included in the raw field.
				assert.NotNil(t, volume.Raw)
				volume.Raw = ""

				assert.True(t, ok)
				return assert.Equal(t, expectedVolume, volume)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := tt.fields.discovery

			got, err := d.List()
			tt.wantErr(t, err)

			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}
