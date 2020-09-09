package k8s

import (
	v1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobOpts struct {
	Image        string
	BackoffLimit int
	Restart      core.RestartPolicy
	Env          map[string]string
	Labels       map[string]string
}

func NewJob(name, namespace string, opts JobOpts) *v1.Job {
	backoffLimit := int32(opts.BackoffLimit)
	envVar := []core.EnvVar{}
	terminationSecs := int64(60)

	// add env vars if any provided
	for name, val := range opts.Env {
		envVar = append(envVar, core.EnvVar{
			Name:  name,
			Value: val,
		})
	}

	job := &v1.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    opts.Labels,
		},
		Spec: v1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Name:      name,
					Namespace: namespace,
					Labels:    opts.Labels,
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{Name: name, Image: opts.Image, Env: envVar},
					},
					RestartPolicy:                 opts.Restart,
					TerminationGracePeriodSeconds: &terminationSecs,
				},
			},
		},
	}

	return job
}
