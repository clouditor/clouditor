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

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/testdata"
	"clouditor.io/clouditor/v2/internal/testutil/assert"

	"google.golang.org/protobuf/testing/protocmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewKubernetesComputeDiscovery(t *testing.T) {
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
			want: &k8sComputeDiscovery{
				k8sDiscovery: k8sDiscovery{},
			},
		},
		{
			name: "Happy path",
			args: args{
				intf:                  &fake.Clientset{},
				CertificationTargetID: testdata.MockCertificationTargetID1,
			},
			want: &k8sComputeDiscovery{
				k8sDiscovery: k8sDiscovery{
					intf: &fake.Clientset{},
					csID: testdata.MockCertificationTargetID1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewKubernetesComputeDiscovery(tt.args.intf, tt.args.CertificationTargetID)
			assert.Equal(t, tt.want, got, assert.CompareAllUnexported())
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
		want    assert.Want[[]ontology.IsResource]
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Happy path",
			fields: fields{
				NewKubernetesComputeDiscovery(client, testdata.MockCertificationTargetID1),
			},
			want: func(t *testing.T, got []ontology.IsResource) bool {
				container, ok := got[0].(*ontology.Container)
				if !assert.True(t, ok) {
					return false
				}
				// Create expected ontology.Container
				expectedContainer := &ontology.Container{
					Id:     podID,
					Name:   podName,
					Labels: podLabel,
					NetworkInterfaceIds: []string{
						podNamespace,
					},
				}

				// We need to ignore creation_time in the comparison because it is random and raw because it includes the creation_time
				assert.NotNil(t, container.CreationTime)
				assert.NotEmpty(t, container.Raw)
				assert.Equal(t, expectedContainer, container, protocmp.IgnoreFields(&ontology.Container{}, "creation_time", "raw"))

				// Check volume
				volume, ok := got[1].(*ontology.BlockStorage)
				assert.True(t, ok)

				// Create expected ontology.BlockStorage
				expectedVolume := &ontology.BlockStorage{
					Id:               volumeName,
					Name:             volumeName,
					CreationTime:     nil,
					AtRestEncryption: &ontology.AtRestEncryption{},
				}

				// We need to ignore raw because it contains the random creation time of the pod
				return assert.NotEmpty(t, container.Raw) && assert.Equal(t, expectedVolume, volume, protocmp.IgnoreFields(&ontology.BlockStorage{}, "raw"))
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

func Test_k8sComputeDiscovery_handlePodVolume(t *testing.T) {
	type fields struct {
		k8sDiscovery k8sDiscovery
	}
	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []ontology.IsResource
	}{
		{
			name:   "file storage",
			fields: fields{},
			args: args{
				pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{
							{
								Name: "test",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/tmp",
									},
								},
							},
						},
					},
				},
			},
			want: []ontology.IsResource{
				&ontology.FileStorage{
					Id:               "test",
					Name:             "test",
					AtRestEncryption: &ontology.AtRestEncryption{},
					Raw:              `{"*v1.Pod":[{"metadata":{"creationTimestamp":null},"spec":{"volumes":[{"name":"test","hostPath":{"path":"/tmp"}}],"containers":null},"status":{}}],"*v1.Volume":[{"name":"test","hostPath":{"path":"/tmp"}}]}`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &k8sComputeDiscovery{
				k8sDiscovery: tt.fields.k8sDiscovery,
			}

			got := d.handlePodVolume(tt.args.pod)
			assert.Equal(t, tt.want, got)
		})
	}
}
