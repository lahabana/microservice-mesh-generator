package main

import (
	"flag"
	"fmt"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s/apiplay"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s/fakeservice"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/yaml"
	"github.com/lahabana/microservice-mesh-generator/pkg/random"
	"github.com/lahabana/microservice-mesh-generator/pkg/version"
	"os"
	"time"
)

func main() {
	numServices := flag.Int("numServices", 20, "The number of services to use")
	minReplicas := flag.Int("minReplicas", 1, "The minimum number of replicas to use (will pick a number between min and max)")
	maxReplicas := flag.Int("maxReplicas", 1, "The max number of replicas to use (will pick a number between min and max)")
	percentEdge := flag.Int("percentEdge", 50, "The for an edge between 2 nodes to exist (100 == sure)")
	seed := flag.Int64("seed", time.Now().Unix(), "the seed for the random generate (set to now by default)")
	k8sNamespace := flag.String("k8sNamespace", "microservice-mesh", "The namespace to use (only useful if output is `k8s`)")
	k8sApp := flag.String("k8sApp", "api-play", "The app to use can be api-play or fake-service (only useful if output is `k8s`)")
	output := flag.String("output", "yaml", "output format (k8s,dot,mermaid,yaml,json)")
	flag.Parse()

	commentMarker := "#"
	var generator apis.Generator
	switch *output {
	case "k8s":
		var opts []k8s.Option
		switch *k8sApp {
		case "api-play":
			opts = apiplay.GeneratorOpts()
		case "fake-service":
			opts = fakeservice.GeneratorOpts()
		default:
			panic(fmt.Errorf("invalid k8sApp '%s' supported: api-play or fake-service", *k8sApp))
		}
		opts = append(opts, k8s.WithNamespace(*k8sNamespace))
		k8sGenerator, err := k8s.NewGenerator(opts...)
		if err != nil {
			panic(err)
		}
		generator = k8sGenerator
	case "dot":
		generator = apis.DotGenerator
	case "mermaid":
		commentMarker = "%%"
		generator = apis.MermaidGenerator
	case "yaml":
		generator = yaml.Generator
	case "json":
		generator = apis.JsonGenerator
	default:
		panic(fmt.Errorf("format '%s' not supported accepted format: k8s, yaml, dot, mermaid, json", *output))
	}
	serviceGraph := random.GenerateRandomMesh(*seed, *numServices, *percentEdge, *minReplicas, *maxReplicas)
	fmt.Printf("%s name:'%s',version:'%s',commit:'%s',seed: %d\n", commentMarker, version.Name, version.Version, version.Commit, *seed)
	err := generator.Apply(os.Stdout, serviceGraph)
	if err != nil {
		panic(any(err))
	}
}
