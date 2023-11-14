package fakeservice

import (
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s"
	v1 "k8s.io/api/core/v1"
	"strings"
)

func GeneratorOpts() []k8s.Option {
	return []k8s.Option{
		k8s.WithPort(9090),
		k8s.WithFormatters(k8s.SimpleFormatters("fake-service")),
		k8s.WithImage("nicholasjackson/fake-service:v0.26.0"),
		k8s.WithPodTemplateSpecMutator(mutatePodTemplate),
	}
}

func mutatePodTemplate(formatters k8s.Formatters, svc apis.Service, template *v1.PodTemplateSpec) error {
	var uris []string
	for _, v := range svc.Edges {
		uris = append(uris, formatters.Url(v, 9090))
	}
	template.Spec.Containers[0].Env = append(template.Spec.Containers[0].Env,
		v1.EnvVar{
			Name:  "SERVICE",
			Value: formatters.Name(svc.Idx),
		},
		v1.EnvVar{
			Name:  "UPSTREAM_URIS",
			Value: strings.Join(uris, ","),
		},
	)
	return nil
}
