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

func (k k8sComputeDiscovery) List() ([]voc.IsResource, error) {
	var list []voc.IsResource

	pods, err := k.intf.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list ingresses: %v", err)
	}

	for i := range pods.Items {
		c := k.handlePod(&pods.Items[i])

		log.Infof("Adding container %+v", c)

		list = append(list, c)
	}

	return list, nil
}

func (k k8sComputeDiscovery) handlePod(pod *v1.Pod) voc.IsCompute {
	return &voc.ContainerResource{
		ComputeResource: voc.ComputeResource{
			Resource: voc.Resource{
				ID:           fmt.Sprintf("/namespaces/%s/containers/%s", pod.Namespace, pod.Name),
				Name:         pod.Name,
				CreationTime: pod.CreationTimestamp.Unix(),
			}},
	}

}
