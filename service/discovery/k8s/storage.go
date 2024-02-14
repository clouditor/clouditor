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
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/ontology"
	"google.golang.org/protobuf/types/known/timestamppb"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sStorageDiscovery struct{ k8sDiscovery }

func NewKubernetesStorageDiscovery(intf kubernetes.Interface, cloudServiceID string) discovery.Discoverer {
	return &k8sStorageDiscovery{k8sDiscovery{intf, cloudServiceID}}
}

func (*k8sStorageDiscovery) Name() string {
	return "Kubernetes Storage"
}

func (*k8sStorageDiscovery) Description() string {
	return "Discover Kubernetes storage resources."
}

func (d *k8sStorageDiscovery) List() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	// Get persistent volumes
	// Note: Volumes exist in the context of a pod and cannot be created on its own, PersistentVolumes are first class objects with its own lifecycle.
	pvc, err := d.intf.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range pvc.Items {
		p := d.handlePV(&pvc.Items[i])
		if p != nil {
			log.Infof("Adding volume %+v", p)
			list = append(list, p)
		}
	}

	if list == nil {
		log.Debugf("No Kubernetes persistent volumes available")
	}

	return list, nil
}

// handlePVC returns all PersistentVolumes
func (d *k8sStorageDiscovery) handlePV(pv *v1.PersistentVolume) ontology.IsResource {
	vs := pv.Spec.PersistentVolumeSource

	// TODO(all): Define all volume types
	// LocalVolumeSource
	// PersistentVolumeClaimVolumeSource
	// DownwardAPIVolumeSource
	// ConfigMapVolumeSource
	// VsphereVirtualDiskVolumeSource
	// QuobyteVolumeSource
	// PhotonPersistentDiskVolumeSource
	// ProjectedVolumeSource
	// ScaleIOVolumeSource
	// CSIVolumeSource -> CSI was developed as a standard for exposing arbitrary block and file storage storage systems to containerized workloads on Container Orchestration Systems (COs) like Kubernetes. (https://kubernetes.io/blog/2019/01/15/container-storage-interface-ga/)

	// Deprecated:
	// GitRepoVolumeSource is deprecated
	// cinder - Cinder (OpenStack block storage) (deprecated in v1.18)
	// flexVolume - FlexVolume (deprecated in v1.23)
	// flocker - Flocker storage (deprecated in v1.22)
	// quobyte - Quobyte volume (deprecated in v1.22)
	// storageos - StorageOS volume (deprecated in v1.22)
	if vs.AWSElasticBlockStore != nil || vs.AzureDisk != nil || vs.Cinder != nil || vs.FlexVolume != nil || vs.CephFS != nil || vs.Glusterfs != nil || vs.GCEPersistentDisk != nil || vs.RBD != nil || vs.StorageOS != nil || vs.FC != nil || vs.PortworxVolume != nil || vs.ISCSI != nil || vs.Flocker != nil {
		v := &ontology.BlockStorage{
			Id:               string(pv.UID),
			Name:             pv.Name,
			CreationTime:     timestamppb.New(pv.CreationTimestamp.Time),
			Labels:           pv.Labels,
			AtRestEncryption: &ontology.AtRestEncryption{},
			Raw:              discovery.Raw(pv),
		}

		return v
	} else if vs.AzureFile != nil || vs.NFS != nil || vs.HostPath != nil {
		v := &ontology.FileStorage{
			Id:               string(pv.UID),
			Name:             pv.Name,
			CreationTime:     timestamppb.New(pv.CreationTimestamp.Time),
			Labels:           pv.Labels,
			AtRestEncryption: &ontology.AtRestEncryption{},
			Raw:              discovery.Raw(pv),
		}

		return v
	}

	return nil
}
