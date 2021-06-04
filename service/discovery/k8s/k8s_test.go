package k8s_test

//go:generate mockgen -package k8s_test -destination mock_k8s_test.go k8s.io/client-go/kubernetes/typed/core/v1 CoreV1Interface,PodInterface
