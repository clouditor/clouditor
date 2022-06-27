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
		list = append(list, v...)

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

// handleVolume returns all persistens volume claims connected to a pod
func (k8sComputeDiscovery) handleVolume(pod *v1.Pod) []voc.IsCloudResource {

	var (
		volumes []voc.IsCloudResource
	)

	// TODO(all): Do we have to differentiate between between persistend volume claim,persistent volumes and storage classes?
	// TODO(all): The ID, region, label and atRestEncryption information we have to get directly from the related storage, but I think we do not have credentials for the other providers in Clouditor?
	for _, vol := range pod.Spec.Volumes {
		s := &voc.Storage{
			Resource: &voc.Resource{
				ID:           voc.ResourceID(vol.Name), // The ID we have to get directly from the related storage
				Name:         vol.Name,
				CreationTime: pod.CreationTimestamp.Unix(), // The CreationTime we have to get directly from the related storage
				Type:         []string{"BlockStorage", "Storage", "Resource"},
				Labels:       pod.Labels,
			},
			AtRestEncryption: &voc.AtRestEncryption{}, // Not able to get the AtRestEncryption information, that must be retrieved directly from the storage
		}

		// TODO(anatheka): Possible to use generics for the follwing if?
		// TODO(all): Define all volume types
		// PersistentVolumeClaimVolumeSource
		// DownwardAPIVolumeSource
		// ConfigMapVolumeSource
		// VsphereVirtualDiskVolumeSource
		// QuobyteVolumeSource
		// PhotonPersistentDiskVolumeSource
		// ProjectedVolumeSource
		// ScaleIOVolumeSource
		// CSIVolumeSource -> CSI was developed as a standard for exposing arbitrary block and file storage storage systems to containerized workloads on Container Orchestration Systems (COs) like Kubernetes. (https://kubernetes.io/blog/2019/01/15/container-storage-interface-ga/)
		// EphemeralVolumeSource

		// Deprecated
		// GitRepoVolumeSource is deprecated
		// cinder - Cinder (OpenStack block storage) (deprecated in v1.18)
		// flexVolume - FlexVolume (deprecated in v1.23)
		// flocker - Flocker storage (deprecated in v1.22)
		// quobyte - Quobyte volume (deprecated in v1.22)
		// storageos - StorageOS volume (deprecated in v1.22)
		if vol.AWSElasticBlockStore != nil || vol.AzureDisk != nil || vol.Cinder != nil || vol.FlexVolume != nil || vol.CephFS != nil || vol.Glusterfs != nil || vol.GCEPersistentDisk != nil || vol.RBD != nil || vol.StorageOS != nil || vol.FC != nil || vol.PortworxVolume != nil || vol.ISCSI != nil || vol.Flocker != nil {
			v := &voc.BlockStorage{
				Storage: s,
			}
			volumes = append(volumes, v)
		} else if vol.AzureFile != nil || vol.EmptyDir != nil || vol.NFS != nil || vol.HostPath != nil || vol.Secret != nil {
			v := &voc.FileStorage{
				Storage: s,
			}
			volumes = append(volumes, v)
		} else {
			continue
		}
	}

	return volumes
}
