package broker

import (
	"fmt"
	"os"
	"strings"

	"github.com/interconnectedcloud/qdr-image/test/k8s/utils"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/constants"
	"github.com/interconnectedcloud/qdr-image/test/k8s/utils/k8s"
	v1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ActiveMQArtemisImageEnvVar = "ACTIVEMQ_ARTEMIS_IMAGE"

type ActiveMQArtemisDeploymentOpts struct {
	k8s.DeploymentOpts
	Name   string
	User   string
	Pass   string
	Queues []string
}

// NewDeployment creates an ActiveMQArtemis Deployment
func NewDeployment(namespace string, opts ActiveMQArtemisDeploymentOpts) (*v1.Deployment, error) {

	var err error

	// Validating mandatory fields
	if utils.StrEmpty(namespace) {
		err := fmt.Errorf("namespace is required")
		return nil, err
	}
	if utils.StrEmpty(opts.Name) {
		err := fmt.Errorf("ActiveMQArtemisDeploymentOpts.Name is required")
		return nil, err
	}

	// Default values
	image := utils.StrDefault(constants.ActiveMQArtemisImage, opts.Image, os.Getenv(ActiveMQArtemisImageEnvVar))
	user := utils.StrDefault("admin", opts.User)
	pass := utils.StrDefault("admin", opts.Pass)

	// Static definitions for ActiveMQ Artemis Deployment
	replicas := int32(1)
	terminationSecs := int64(60)
	extraArgs := []string{
		"--host", "0.0.0.0",
		"--http-host", "0.0.0.0",
		"--allow-anonymous",
	}
	pullPolicy := v13.PullAlways
	restartPolicy := v13.RestartPolicyAlways

	// Parsing extra args
	if len(opts.Queues) > 0 {
		extraArgs = append(extraArgs, "--queues")
		for _, q := range opts.Queues {
			extraArgs = append(extraArgs, q)
		}
	}

	var d *v1.Deployment = &v1.Deployment{
		ObjectMeta: v12.ObjectMeta{
			Name:      opts.Name,
			Namespace: namespace,
			Labels:    opts.Labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &v12.LabelSelector{
				MatchLabels: opts.Labels,
			},
			Template: v13.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Labels: opts.Labels,
				},
				Spec: v13.PodSpec{
					Containers: []v13.Container{
						{Name: opts.Name, Image: image, ImagePullPolicy: pullPolicy,
							Env: []v13.EnvVar{
								{Name: "AMQ_USER", Value: user},
								{Name: "AMQ_PASSWORD", Value: pass},
								{Name: "AMQ_EXTRA_ARSG", Value: strings.Join(extraArgs, " ")},
							}},
					},
					RestartPolicy:                 restartPolicy,
					TerminationGracePeriodSeconds: &terminationSecs,
				},
			},
		},
	}

	return d, err
}
