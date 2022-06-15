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
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sComputeDiscovery struct{ k8sDiscovery }

func NewKubernetesComputeDiscovery(intf kubernetes.Interface) discovery.Discoverer {
	return &k8sComputeDiscovery{k8sDiscovery{intf}}
}

func (*k8sComputeDiscovery) Name() string {
	return "Kubernetes Compute"
}

func (*k8sComputeDiscovery) Description() string {
	return "Discover Kubernetes compute resources."
}

func (d k8sComputeDiscovery) List() ([]voc.IsCloudResource, error) {
	var list []voc.IsCloudResource

	pods, err := d.intf.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range pods.Items {
		c := d.handlePod(&pods.Items[i])
		log.Infof("Adding container %+v", c)
		list = append(list, c)

		v := d.handleVolume(&pods.Items[i])
		log.Infof("Adding volume %+v", v)
		list = append(list, v)

	}

	return list, nil
}

func (k8sComputeDiscovery) handlePod(pod *v1.Pod) *voc.Container {
	r := &voc.Container{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(getContainerResourceID(pod)),
				Name:         pod.Name,
				CreationTime: pod.CreationTimestamp.Unix(),
				Type:         []string{"Container", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: "", // TODO(all) Add region to k8s container
				},
			}},
	}

	r.NetworkInterface = append(r.NetworkInterface, voc.ResourceID(pod.Namespace))

	return r

}

func getContainerResourceID(pod *v1.Pod) string {
	return fmt.Sprintf("/namespaces/%s/containers/%s", pod.Namespace, pod.Name)
}

// handleVolume return all persistens volumes connected to a pod
func (k8sComputeDiscovery) handleVolume(pod *v1.Pod) *voc.BlockStorage {

	var name string

	for i := range pod.Spec.Volumes {
		if pod.Spec.Volumes[i].PersistentVolumeClaim != nil {
			name = pod.Spec.Volumes[i].PersistentVolumeClaim.ClaimName
		}
	}

	// TODO(anatheka): Can we get all volumes instead of the connected persistent volume claims
	// Difference between persistend volume claim,persistent volumes and storage classes
	s := &voc.BlockStorage{
		Storage: &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(getContainerResourceID(pod)), //Fix
				Name:         name,
				CreationTime: pod.CreationTimestamp.Unix(), // Fix
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: "", // TODO(all) Add region to k8s volume
				},
			},
			AtRestEncryption: &voc.AtRestEncryption{},
		},
	}

	return s
}
