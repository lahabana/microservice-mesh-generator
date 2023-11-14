package yaml

import (
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"gopkg.in/yaml.v3"
	"io"
)

// Generator outputs the service graph as a yaml.
var Generator = apis.GeneratorFunc(func(writer io.Writer, svc apis.ServiceGraph) error {
	return yaml.NewEncoder(writer).Encode(svc)
})
