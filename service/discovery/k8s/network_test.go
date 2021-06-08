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

package k8s_test

import (
	"context"
	"testing"

	"clouditor.io/clouditor/service/discovery/k8s"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListIngresses(t *testing.T) {
	client := fake.NewSimpleClientset()

	_, err := client.NetworkingV1().Ingresses("my-namespace").Create(context.TODO(), &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: "my-ingress", CreationTimestamp: metav1.Now()},
		Spec: networkingv1.IngressSpec{Rules: []networkingv1.IngressRule{{
			Host: "myhost",
			IngressRuleValue: networkingv1.IngressRuleValue{HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{{
					Path: "test",
				}},
			},
			},
		}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	_, err = client.NetworkingV1().Ingresses("my-namespace").Create(context.TODO(), &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: "my-other-ingress", CreationTimestamp: metav1.Now()},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: "myhost",
				IngressRuleValue: networkingv1.IngressRuleValue{HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{{
						Path: "test",
					}}},
				},
			}},
			TLS: []networkingv1.IngressTLS{{Hosts: []string{"myhost"}}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting ingress add: %v", err)
	}

	_, err = client.CoreV1().Services("my-namespace").Create(context.TODO(), &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "my-service", CreationTimestamp: metav1.Now()},
		Spec: v1.ServiceSpec{
			ClusterIPs: []string{"127.0.0.1"},
			Ports:      []v1.ServicePort{{Port: 80}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting service add: %v", err)
	}

	d := k8s.NewKubernetesNetworkDiscovery(client)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	service, ok := list[0].(*voc.NetworkService)

	assert.True(t, ok)
	assert.Equal(t, "my-service", service.Name)
	assert.Equal(t, "/namespaces/my-namespace/services/my-service", service.ID)
	assert.Equal(t, []int16{80}, service.Ports)
	assert.Equal(t, []string{"127.0.0.1"}, service.IPs)

	lb, ok := list[1].(*voc.LoadBalancerResource)

	assert.True(t, ok)
	assert.Equal(t, "my-ingress", lb.Name)
	assert.Equal(t, "/namespaces/my-namespace/ingresses/my-ingress", lb.ID)
	assert.Equal(t, "http://myhost/test", lb.HttpEndpoints[0].ID)

	lb, ok = list[2].(*voc.LoadBalancerResource)

	assert.True(t, ok)
	assert.Equal(t, "my-other-ingress", lb.Name)
	assert.Equal(t, "/namespaces/my-namespace/ingresses/my-other-ingress", lb.ID)
	assert.Equal(t, "https://myhost/test", lb.HttpEndpoints[0].ID)
	assert.NotNil(t, lb.HttpEndpoints[0].TransportEncryption)
}
