package apiplay

import (
	"encoding/json"
	"fmt"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s"
	v1 "k8s.io/api/core/v1"
)

func GeneratorOpts() []k8s.Option {
	return []k8s.Option{
		k8s.WithPort(8080),
		k8s.WithFormatters(k8s.SimpleFormatters("api-play")),
		k8s.WithImage("ghcr.io/lahabana/api-play:main"),
		k8s.WithConfigMapGenerator(configMapGenerator(8080)),
		k8s.WithPodTemplateSpecMutator(podTemplateMutator),
	}
}

func podTemplateMutator(f k8s.Formatters, svc apis.Service, template *v1.PodTemplateSpec) error {
	template.Spec.Containers[0].Args = []string{"-config-file", "/etc/config/config.yaml"}
	return nil
}

func configMapGenerator(port int) func(formatters k8s.Formatters, svc apis.Service) (string, error) {
	return func(formatters k8s.Formatters, svc apis.Service) (string, error) {
		calls := []map[string]string{}
		for _, s := range svc.Edges {
			calls = append(calls, map[string]string{
				"url": formatters.Url(s, port) + "/api/dynamic/microservice_mesh",
			})

		}
		res, err := json.MarshalIndent(map[string][]map[string]interface{}{
			"apis": {
				{
					"path": "microservice_mesh",
					"conf": map[string]interface{}{
						"body": fmt.Sprintf("I am Service %d and I have %d replicas", svc.Idx, svc.Replicas),
						"call": calls,
					},
				},
			},
		}, "", "  ")
		return string(res), err
	}
}
