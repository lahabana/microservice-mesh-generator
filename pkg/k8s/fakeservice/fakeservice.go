package fakeservice

import (
	"fmt"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/k8s"
	v1 "k8s.io/api/core/v1"
	"strings"
)

const baseName = "fake-svc"

func Encoder(opts ...k8s.Option) (k8s.Encoder, error) {
	opts = append([]k8s.Option{
		k8s.WithPort(9090),
		k8s.WithBaseName(baseName),
		k8s.WithImage("nicholasjackson/fake-service:v0.26.0"),
		k8s.WithPodTemplateSpecMutator(mutatePodTemplate),
	}, opts...)
	enc, err := k8s.GenericEncoder(opts...)
	if err != nil {
		return enc, err
	}
	return enc, nil
}

func mutatePodTemplate(name string, svc apis.Service, template *v1.PodTemplateSpec) error {
	var uris []string
	for _, v := range svc.Edges {
		uris = append(uris, fmt.Sprintf("http://%s-%03d:%d", baseName, v, 9090))
	}
	template.Spec.Containers[0].Args = nil
	template.Spec.Containers[0].Env = append(template.Spec.Containers[0].Env,
		v1.EnvVar{
			Name:  "SERVICE",
			Value: name,
		},
		v1.EnvVar{
			Name:  "UPSTREAM_URIS",
			Value: strings.Join(uris, ","),
		},
	)
	return nil
}
