package k8s

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewServiceClusterIP create a Service instance of a ClusterIP type
// using the given labels and selectors.
func NewServiceClusterIP(name, namespace string, ports []int, labels, selectorLabels map[string]string) *v1.Service {
	// Creating a simple Service with ClusterIP type
	s := &v1.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Selector: selectorLabels,
			Type:     v1.ServiceTypeClusterIP,
			Ports:    []v1.ServicePort{},
		},
	}

	// Adding ports
	for _, port := range ports {
		s.Spec.Ports = append(s.Spec.Ports, v1.ServicePort{
			Name: strconv.Itoa(port),
			Port: int32(port),
		})
	}

	return s
}
