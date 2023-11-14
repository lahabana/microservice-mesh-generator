package main

import (
	"context"
	"flag"
	"github.com/lahabana/microservice-mesh-generator/internal/generate"
	"github.com/lahabana/microservice-mesh-generator/internal/restapi"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0 -config openapi.cfg.yaml openapi.yaml

func main() {
	config := generate.DefaultConfig()
	numServices := flag.Int("numServices", 5, "The number of services to use")
	minReplicas := flag.Int("minReplicas", 2, "The minimum number of replicas to use (will pick a number between min and max)")
	maxReplicas := flag.Int("maxReplicas", 2, "The max number of replicas to use (will pick a number between min and max)")
	percentEdge := flag.Int("percentEdge", 50, "The for an edge between 2 nodes to exist (100 == sure)")
	flag.Int64Var(&config.Seed, "seed", config.Seed, "the seed for the random generate (set to now by default)")
	flag.StringVar(&config.K8sNamespace, "k8sNamespace", config.K8sNamespace, "The namespace to use (only useful if output is `k8s`)")
	flag.StringVar(&config.K8sApp, "k8sApp", config.K8sApp, "The app to use can be api-play or fake-service (only useful if output is `k8s`)")
	flag.StringVar(&config.Output, "output", config.Output, "output format (k8s,dot,mermaid,yaml,json)")
	asServer := flag.Bool("server", false, "whether to run this tool as a hosted server")
	flag.Parse()

	if *asServer {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := restapi.Start(ctx)
		if err != nil {
			panic(err)
		}
		return
	}
	err := generate.Run(config, func(seed int64) (apis.ServiceGraph, error) {
		mesh := apis.GenerateRandomMesh(seed, *numServices, *percentEdge, *minReplicas, *maxReplicas)
		return mesh, nil
	})
	if err != nil {
		panic(any(err))
	}
}
