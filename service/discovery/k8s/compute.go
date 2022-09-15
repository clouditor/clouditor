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

	// Get pods
	pods, err := d.intf.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range pods.Items {
		// Get virtual machines
		c := d.handlePod(&pods.Items[i])
		log.Infof("Adding container %+v", c)
		list = append(list, c)

		// Get all volumes conntected to the specific pod
		v := d.handlePodVolume(&pods.Items[i])

		if len(v) != 0 {
			log.Infof("Adding pod volume %+v", v)
			list = append(list, v...)
		}
	}

	return list, nil
}

// handlePod returns all existing pods
func (k8sComputeDiscovery) handlePod(pod *v1.Pod) *voc.Container {
	r := &voc.Container{
		Compute: &voc.Compute{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(getContainerResourceID(pod)),
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         pod.Name,
				CreationTime: pod.CreationTimestamp.Unix(),
				Type:         []string{"Container", "Compute", "Resource"},
				GeoLocation: voc.GeoLocation{
					Region: "", // TODO(all) Add region to k8s container
				},
				Labels: pod.Labels,
			},
		},
	}

	r.NetworkInterface = append(r.NetworkInterface, voc.ResourceID(pod.Namespace))

	return r

}

func getContainerResourceID(pod *v1.Pod) string {
	return fmt.Sprintf("/namespaces/%s/containers/%s", pod.Namespace, pod.Name)
}

// handleVolume returns all volumes connected to a pod
func (k8sComputeDiscovery) handlePodVolume(pod *v1.Pod) []voc.IsCloudResource {

	var (
		volumes []voc.IsCloudResource
	)

	// TODO(all): Do we have to differentiate between between persistend volume claim,persistent volumes and storage classes?
	// TODO(all): The ID, region, label and atRestEncryption information we have to get directly from the related storage, but I think we do not have credentials for the other providers in Clouditor?
	for _, vol := range pod.Spec.Volumes {
		s := &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(vol.Name), // The ID we have to get directly from the related storage
				ServiceID:    discovery.DefaultCloudServiceID,
				Name:         vol.Name,
				CreationTime: 0, // The CreationTime we have to get directly from the related storage
				Type:         []string{"BlockStorage", "Storage", "Resource"},
			},
			AtRestEncryption: &voc.AtRestEncryption{}, // Not able to get the AtRestEncryption information, that must be retrieved directly from the storage
		}

		v := addVolumeSource(s, vol.VolumeSource)
		if v != nil {
			volumes = append(volumes, v)
		}
	}

	return volumes
}
