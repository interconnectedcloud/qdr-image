package k8s

import (
	"fmt"

	"github.com/interconnectedcloud/qdr-image/test/k8s/utils"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentOpts struct {
	Image         string
	Labels        map[string]string
	RestartPolicy v13.RestartPolicy
	Command       []string
	Args          []string
	EnvVars       []v13.EnvVar
}

func NewDeployment(name, namespace string, opts DeploymentOpts) (*v1.Deployment, error) {

	var err error

	// Validating mandatory fields
	if utils.StrEmpty(namespace) {
		err := fmt.Errorf("namespace is required")
		return nil, err
	}
	if utils.StrEmpty(name) {
		err := fmt.Errorf("deployment name is required")
		return nil, err
	}
	if utils.StrEmpty(opts.Image) {
		err := fmt.Errorf("image is required")
		return nil, err
	}

	// Container to use
	containers := []v13.Container{
		{Name: name, Image: opts.Image, ImagePullPolicy: v13.PullAlways, Env: opts.EnvVars},
	}
	// Customize commands and arguments if any informed
	if len(opts.Command) > 0 {
		containers[0].Command = opts.Command
	}
	if len(opts.Args) > 0 {
		containers[0].Args = opts.Args
	}

	d := &v1.Deployment{
		ObjectMeta: v12.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    opts.Labels,
		},
		Spec: v1.DeploymentSpec{
			Selector: &v12.LabelSelector{
				MatchLabels: opts.Labels,
			},
			Template: v13.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels: opts.Labels,
				},
				Spec: v13.PodSpec{
					Containers:    containers,
					RestartPolicy: opts.RestartPolicy,
				},
			},
		},
	}

	return d, err
}
