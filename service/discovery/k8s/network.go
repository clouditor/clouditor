package k8s

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sNetworkDiscovery struct{ k8sDiscovery }

func NewKubernetesNetworkDiscovery(intf kubernetes.Interface) discovery.Discoverer {
	return &k8sNetworkDiscovery{k8sDiscovery{intf}}
}

func (d *k8sNetworkDiscovery) Name() string {
	return "Kubernetes Network"
}

func (d *k8sNetworkDiscovery) Description() string {
	return "Discover Kubernetes network resources."
}

func (k k8sNetworkDiscovery) List() ([]voc.IsResource, error) {
	var list []voc.IsResource

	services, err := k.intf.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list services: %v", err)
	}

	for i := range services.Items {
		c := k.handleService(&services.Items[i])

		log.Infof("Adding service %+v", c)

		list = append(list, c)
	}

	ingresses, err := k.intf.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range ingresses.Items {
		c := k.handleIngress(&ingresses.Items[i])

		log.Infof("Adding ingress %+v", c)

		list = append(list, c)
	}

	return list, nil
}

func (k k8sNetworkDiscovery) handleService(service *corev1.Service) voc.IsCompute {
	var ports []int16

	for _, v := range service.Spec.Ports {
		ports = append(ports, int16(v.Port))
	}

	return &voc.NetworkService{
		Resource: voc.Resource{
			ID:           fmt.Sprintf("/namespaces/%s/services/%s", service.Namespace, service.Name),
			Name:         service.Name,
			CreationTime: service.CreationTimestamp.Unix(),
		},
		IPs:   service.Spec.ClusterIPs,
		Ports: ports,
	}
}

func (k k8sNetworkDiscovery) handleIngress(ingress *v1.Ingress) voc.IsCompute {
	var url = fmt.Sprintf("%s/%s", ingress.Spec.Rules[0].Host, ingress.Spec.Rules[0].HTTP.Paths[0].Path)
	var te *voc.TransportEncryption

	if ingress.Spec.TLS == nil {
		url = fmt.Sprintf("http://%s", url)
	} else {
		url = fmt.Sprintf("https://%s", url)

		te = &voc.TransportEncryption{
			Enforced:   true,
			Encryption: voc.Encryption{Enabled: true},
		}
	}

	return &voc.HttpEndpoint{
		Resource: voc.Resource{
			ID:           url,
			Name:         ingress.Name,
			CreationTime: ingress.CreationTimestamp.Unix(),
		},
		URL:                 url,
		TransportEncryption: te,
	}
}
