package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Service struct {
	Idx      int   `yaml:"idx" json:"idx"`
	Edges    []int `yaml:"edges" json:"edges"`
	Replicas int   `yaml:"replicas" json:"replicas"`
}

type ServiceGraph struct {
	Services         []Service `yaml:"services" json:"services"`
	GenerationParams string    `yaml:"generationParams" json:"generationParams"`
}

// Generator generates the graph is a custom format
type Generator interface {
	Apply(writer io.Writer, svc ServiceGraph) error
}

// GeneratorFunc a shorthand for a regular generator
type GeneratorFunc func(writer io.Writer, svc ServiceGraph) error

func (f GeneratorFunc) Apply(writer io.Writer, svc ServiceGraph) error {
	return f(writer, svc)
}

// DotGenerator outputs the service graph in dot format.
var DotGenerator = GeneratorFunc(func(writer io.Writer, s ServiceGraph) error {
	var allEdges []string
	for _, srv := range s.Services {
		for _, other := range srv.Edges {
			allEdges = append(allEdges, fmt.Sprintf("%d -> %d;", srv.Idx, other))
		}
	}
	_, err := fmt.Fprintf(writer, "digraph{\n%s\n}\n", strings.Join(allEdges, "\n"))
	return err
})

// MermaidGenerator outputs the service graph in mermaid format (this is rendered in github flavoured Markdown).
var MermaidGenerator = GeneratorFunc(func(writer io.Writer, s ServiceGraph) error {
	var allEdges []string
	for _, srv := range s.Services {
		for _, other := range srv.Edges {
			allEdges = append(allEdges, fmt.Sprintf("\t%d --> %d;", srv.Idx, other))
		}
	}
	_, err := fmt.Fprintf(writer, "graph TD;\n%s\n\n", strings.Join(allEdges, "\n"))
	return err
})

// JsonGenerator outputs the service graph in json
var JsonGenerator = GeneratorFunc(func(writer io.Writer, svc ServiceGraph) error {
	return json.NewEncoder(writer).Encode(svc)
})
