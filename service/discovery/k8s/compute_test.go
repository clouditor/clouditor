package k8s_test

import (
	"context"
	"testing"

	"clouditor.io/clouditor/service/discovery/k8s"
	"clouditor.io/clouditor/voc"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListPods(t *testing.T) {
	client := fake.NewSimpleClientset()

	p := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "my-pod", CreationTimestamp: metav1.Now()}}
	_, err := client.CoreV1().Pods("my-namespace").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	d := k8s.NewKubernetesComputeDiscovery(client)

	list, err := d.List()

	assert.Nil(t, err)
	assert.NotNil(t, list)

	container, ok := list[0].(*voc.ContainerResource)

	assert.True(t, ok)
	assert.Equal(t, "my-pod", container.Name)
	assert.Equal(t, "/namespaces/my-namespace/containers/my-pod", container.ID)
}
