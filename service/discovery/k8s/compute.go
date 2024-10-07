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
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sComputeDiscovery struct{ k8sDiscovery }

func NewKubernetesComputeDiscovery(intf kubernetes.Interface, CertificationTargetID string) discovery.Discoverer {
	return &k8sComputeDiscovery{k8sDiscovery{intf, CertificationTargetID}}
}

func (*k8sComputeDiscovery) Name() string {
	return "Kubernetes Compute"
}

func (*k8sComputeDiscovery) Description() string {
	return "Discover Kubernetes compute resources."
}

func (d *k8sComputeDiscovery) List() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// Get pods
	pods, err := d.intf.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range pods.Items {
		// Get virtual machines
		c := d.handlePod(&pods.Items[i])
		log.Infof("Adding container %+v", c.GetId())
		list = append(list, c)

		// Get all volumes conntected to the specific pod
		v := d.handlePodVolume(&pods.Items[i])

		if len(v) != 0 {
			log.Infof("Adding pod volumes %+v", strings.Join(ontology.ResourceIDs(v), ","))
			list = append(list, v...)
		}
	}

	return list, nil
}

// handlePod returns all existing pods
func (d *k8sComputeDiscovery) handlePod(pod *v1.Pod) *ontology.Container {
	r := &ontology.Container{
		Id:           getContainerResourceID(pod),
		Name:         pod.Name,
		CreationTime: timestamppb.New(pod.CreationTimestamp.Time),
		Labels:       pod.Labels,
		Raw:          discovery.Raw(pod),
	}

	r.NetworkInterfaceIds = append(r.NetworkInterfaceIds, pod.Namespace)

	return r
}

func getContainerResourceID(pod *v1.Pod) string {
	return fmt.Sprintf("/namespaces/%s/containers/%s", pod.Namespace, pod.Name)
}

// handleVolume returns all volumes connected to a pod
func (d *k8sComputeDiscovery) handlePodVolume(pod *v1.Pod) []ontology.IsResource {
	var (
		volumes []ontology.IsResource
		v       ontology.IsResource
	)

	// TODO(all): Do we have to differentiate between between persisted volume claim,persistent volumes and storage
	// classes?
	// TODO(all): The ID, region, label and atRestEncryption information we have to get directly from the related
	// storage, but I think we do not have credentials for the other providers in Clouditor?
	for _, vol := range pod.Spec.Volumes {
		vs := vol.VolumeSource

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

		// Deprecated:
		// GitRepoVolumeSource is deprecated
		// cinder - Cinder (OpenStack block storage) (deprecated in v1.18)
		// flexVolume - FlexVolume (deprecated in v1.23)
		// flocker - Flocker storage (deprecated in v1.22)
		// quobyte - Quobyte volume (deprecated in v1.22)
		// storageos - StorageOS volume (deprecated in v1.22)
		if vs.AWSElasticBlockStore != nil || vs.AzureDisk != nil || vs.Cinder != nil || vs.FlexVolume != nil || vs.CephFS != nil || vs.Glusterfs != nil || vs.GCEPersistentDisk != nil || vs.RBD != nil || vs.StorageOS != nil || vs.FC != nil || vs.PortworxVolume != nil || vs.ISCSI != nil || vs.Flocker != nil {
			v = &ontology.BlockStorage{
				Id:           vol.Name, // The ID we have to get directly from the related storage
				Name:         vol.Name,
				CreationTime: nil, // The CreationTime we have to get directly from the related storage
				// anatheka: As I understand it, there are no labels for the volume here, we have to get that from the
				// related storage directly. But we could take the pod labels to which the volume is assigned. I think that
				// makes more sense.
				Labels: nil,
				// Not able to get the AtRestEncryption information, that must be retrieved directly from the storage
				AtRestEncryption: &ontology.AtRestEncryption{},
				Raw:              discovery.Raw(pod, &vol),
			}
		} else if vs.AzureFile != nil || vs.EmptyDir != nil || vs.NFS != nil || vs.HostPath != nil || vs.Secret != nil {
			v = &ontology.FileStorage{
				Id:           vol.Name, // The ID we have to get directly from the related storage
				Name:         vol.Name,
				CreationTime: nil, // The CreationTime we have to get directly from the related storage
				// anatheka: As I understand it, there are no labels for the volume here, we have to get that from the
				// related storage directly. But we could take the pod labels to which the volume is assigned. I think that
				// makes more sense.
				Labels: nil,
				// Not able to get the AtRestEncryption information, that must be retrieved directly from the storage
				AtRestEncryption: &ontology.AtRestEncryption{},
				Raw:              discovery.Raw(pod, &vol),
			}
		}

		if v != nil {
			volumes = append(volumes, v)
		}
	}

	return volumes
}
