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
	"clouditor.io/clouditor/api/ontology"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sNetworkDiscovery struct{ k8sDiscovery }

func NewKubernetesNetworkDiscovery(intf kubernetes.Interface, cloudServiceID string) discovery.Discoverer {
	return &k8sNetworkDiscovery{k8sDiscovery{intf, cloudServiceID}}
}

func (*k8sNetworkDiscovery) Name() string {
	return "Kubernetes Network"
}

func (*k8sNetworkDiscovery) Description() string {
	return "Discover Kubernetes network resources."
}

func (d *k8sNetworkDiscovery) List() ([]ontology.IsResource, error) {
	var list []ontology.IsResource

	services, err := d.intf.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list services: %w", err)
	}

	for i := range services.Items {
		c := d.handleService(&services.Items[i])

		log.Infof("Adding service %+v", c)

		list = append(list, c)
	}

	// TODO Does not get ingresses
	ingresses, err := d.intf.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return list, fmt.Errorf("could not list ingresses: %w", err)
	}

	for i := range ingresses.Items {
		c := d.handleIngress(&ingresses.Items[i])

		log.Infof("Adding ingress %+v", c)

		list = append(list, c)
	}

	return list, nil
}

func (d *k8sNetworkDiscovery) handleService(service *corev1.Service) ontology.IsResource {
	var (
		ports []uint32
	)

	for _, v := range service.Spec.Ports {
		ports = append(ports, uint32(v.Port))
	}

	return &ontology.GenericNetworkService{
		Id:           getNetworkServiceResourceID(service),
		Name:         service.Name,
		CreationTime: timestamppb.New(service.CreationTimestamp.Time),
		Labels:       service.Labels,
		Raw:          discovery.Raw(service),
		Ips:          service.Spec.ClusterIPs,
		Ports:        ports,
	}
}

func getNetworkServiceResourceID(service *corev1.Service) string {
	return fmt.Sprintf("/namespaces/%s/services/%s", service.Namespace, service.Name)
}

func (d *k8sNetworkDiscovery) handleIngress(ingress *v1.Ingress) ontology.IsResource {
	lb := &ontology.LoadBalancer{
		Id:           getLoadBalancerResourceID(ingress),
		Name:         ingress.Name,
		CreationTime: timestamppb.New(ingress.CreationTimestamp.Time),
		Labels:       ingress.Labels,
		Raw:          discovery.Raw(ingress),
		Ports:        []uint32{80, 443},
	}

	for _, rule := range ingress.Spec.Rules {
		lb.Ips = append(lb.Ips, rule.Host)

		for _, path := range rule.HTTP.Paths {
			var url = fmt.Sprintf("%s/%s", rule.Host, path.Path)
			var te *ontology.TransportEncryption

			if ingress.Spec.TLS == nil {
				url = fmt.Sprintf("http://%s", url)
			} else {
				url = fmt.Sprintf("https://%s", url)

				te = &ontology.TransportEncryption{
					Enforced: true,
					Enabled:  true,
				}
			}

			http := &ontology.HttpEndpoint{
				Url:                 url,
				TransportEncryption: te,
			}

			lb.HttpEndpoints = append(lb.HttpEndpoints, http)
		}
	}

	return lb
}

func getLoadBalancerResourceID(ingress *v1.Ingress) string {
	return fmt.Sprintf("/namespaces/%s/ingresses/%s", ingress.Namespace, ingress.Name)
}
