package fakeservice_test

import (
	"bytes"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/k8s"
	"github.com/lahabana/microservice-mesh-generator/pkg/k8s/fakeservice"
	"testing"
)

func TestSimple(t *testing.T) {
	encoder, err := fakeservice.Encoder(k8s.WithNamespace("bar"))
	if err != nil {
		t.Error("failed", err)
	}
	buf := bytes.NewBuffer([]byte{})
	err = encoder.Encode(buf, apis.ServiceGraph{
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
