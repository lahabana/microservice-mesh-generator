package apiplay_test

import (
	"bytes"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s"
	"github.com/lahabana/microservice-mesh-generator/pkg/generators/k8s/apiplay"
	"testing"
)

func TestSimple(t *testing.T) {
	opts := apiplay.GeneratorOpts()
	opts = append(opts, k8s.WithNamespace("foo"))
	encoder, err := k8s.NewGenerator(opts...)
	if err != nil {
		t.Error("failed", err)
	}
	buf := bytes.NewBuffer([]byte{})
	err = encoder.Apply(buf, apis.ServiceGraph{
		Services: []apis.Service{
			{Replicas: 2, Edges: []int{1, 2}, Idx: 0},
			{Replicas: 2, Edges: []int{2}, Idx: 1},
			{Replicas: 2, Edges: []int{3}, Idx: 2},
			{Replicas: 2, Edges: []int{}, Idx: 3},
		},
	})
	if err != nil {
		t.Error("failed", err)
	}
	println(buf.String())
}
