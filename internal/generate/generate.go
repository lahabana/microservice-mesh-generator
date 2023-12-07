package generate

import (
	"fmt"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s/apiplay"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s/fakeservice"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/yaml"
	"github.com/lahabana/microservice-mesh-generator/pkg/version"
	"io"
	"os"
	"time"
)

type Config struct {
	K8s          bool
	Output       string
	K8sApp       string
	K8sNamespace string
	Kuma         bool
	Seed         int64
	Writer       io.Writer
}

var DefaultConfig = func() Config {
	return Config{
		Writer:       os.Stdout,
		Seed:         time.Now().Unix(),
		K8sApp:       "api-play",
		K8sNamespace: "microservice-mesh",
		Output:       "yaml",
		K8s:          false,
	}
}

type InvalidConfError struct {
	msg string
}

func (e *InvalidConfError) Error() string {
	return fmt.Sprintf("invalid param: %s", e.msg)
}

func (e *InvalidConfError) Is(target error) bool {
	_, ok := target.(*InvalidConfError)
	return ok
}

func Run(conf Config, genFn func(seed int64) (apis.ServiceGraph, error)) error {
	commentMarker := "#"
	var generator apis.Generator
	switch conf.Output {
	case "k8s":
		var opts []k8s.Option
		switch conf.K8sApp {
		case "api-play":
			opts = apiplay.GeneratorOpts()
		case "fake-service":
			opts = fakeservice.GeneratorOpts()
		default:
			return &InvalidConfError{msg: fmt.Sprintf("invalid k8sApp '%s' supported: api-play or fake-service", conf.K8sApp)}
		}
		opts = append(opts, k8s.WithNamespace(conf.K8sNamespace))
		k8sGenerator, err := k8s.NewGenerator(opts...)
		if err != nil {
			return err
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
		commentMarker = ""
		generator = apis.JsonGenerator
	default:
		return &InvalidConfError{msg: fmt.Sprintf("format '%s' not supported accepted format: k8s, yaml, dot, mermaid, json", conf.Output)}
	}
	serviceGraph, err := genFn(conf.Seed)
	if err != nil {
		return err
	}
	if err := serviceGraph.Validate(); err != nil {
		return &InvalidConfError{msg: err.Error()}
	}
	if commentMarker != "" {
		_, _ = fmt.Fprintf(conf.Writer, "%s runParameters=package:%s,version:%s,commit:%s,seed:%d\n", commentMarker, version.Name, version.Version, version.Commit, conf.Seed)
		_, _ = fmt.Fprintf(conf.Writer, "%s generationParameters=%s\n", commentMarker, serviceGraph.GenerationParams)
	}
	return generator.Apply(conf.Writer, serviceGraph)
}
