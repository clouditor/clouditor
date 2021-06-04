package k8s

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type k8sNetworkDiscovery struct{ k8sDiscovery }

func NewKubernetesNetworkDiscovery(intf kubernetes.Interface) discovery.Discoverer {
	return &k8sNetworkDiscovery{k8sDiscovery{intf}}
}

func (k k8sNetworkDiscovery) List() ([]voc.IsResource, error) {
	var list []voc.IsResource

	ingresses, err := k.intf.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, ingress := range ingresses.Items {
		c := k.handleIngress(ingress)

		log.Infof("Adding container %+v", c)

		list = append(list, c)
	}

	return list, nil
}

func (k k8sNetworkDiscovery) handleIngress(ingress v1.Ingress) voc.IsCompute {
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
