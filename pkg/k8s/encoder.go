package k8s

import (
	"bytes"
	"fmt"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"sigs.k8s.io/yaml"
)

type Encoder struct {
	CommonSetup       CommonSetup
	WorkloadGenerator WorkloadGenerator
	Serializer        *json.Serializer
}

var DefaultSerializer = json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Yaml: true, Pretty: true, Strict: true})

func (e *Encoder) Encode(writer io.Writer, svc apis.ServiceGraph) error {
	if e.CommonSetup != nil {
		objs, raw, err := e.CommonSetup.Generate(svc)
		if err != nil {
			return err
		}
		if _, err := writer.Write(raw); err != nil {
			return err
		}
		if err := e.encode(writer, objs...); err != nil {
			return err
		}
	}
	for _, s := range svc.Services {
		objs, raw, err := e.WorkloadGenerator.Generate(s)
		if err != nil {
			return &ServiceGeneratorError{idx: s.Idx, err: err}
		}
		if _, err := writer.Write(raw); err != nil {
			return err
		}
		if err := e.encode(writer, objs...); err != nil {
			return &ServiceGeneratorError{idx: s.Idx, err: err}
		}
	}
	return nil
}

func (e *Encoder) encode(writer io.Writer, inputs ...runtime.Object) error {
	for _, in := range inputs {
		_, err := writer.Write([]byte("---\n"))
		if err != nil {
			return err
		}
		b := bytes.Buffer{}
		if err := e.Serializer.Encode(in, &b); err != nil {
			return err
		}

		// Remove "status" from Kubernetes YAMLs
		obj := map[string]interface{}{}
		if err := yaml.Unmarshal(b.Bytes(), &obj); err != nil {
			return err
		}
		delete(obj, "status")
		delete(obj["metadata"].(map[string]interface{}), "creationTimestamp")
		b2, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		if _, err = writer.Write(b2); err != nil {
			return err
		}
	}
	return nil
}

type CommonSetup interface {
	Generate(svcs apis.ServiceGraph) ([]runtime.Object, []byte, error)
}

type CommonSetupFn func(svcs apis.ServiceGraph) ([]runtime.Object, []byte, error)

func (f CommonSetupFn) Generate(svcs apis.ServiceGraph) ([]runtime.Object, []byte, error) {
	return f(svcs)
}

type WorkloadGenerator interface {
	Generate(svc apis.Service) ([]runtime.Object, []byte, error)
}

type WorkloadGeneratorFn func(svc apis.Service) ([]runtime.Object, []byte, error)

func (f WorkloadGeneratorFn) Generate(svc apis.Service) ([]runtime.Object, []byte, error) {
	return f(svc)
}

type ServiceGeneratorError struct {
	idx int
	err error
}

func (s *ServiceGeneratorError) Unwrap() error {
	return s.err
}

func (s *ServiceGeneratorError) Is(target error) bool {
	_, ok := target.(*ServiceGeneratorError)
	return ok
}

func (s *ServiceGeneratorError) Error() string {
	return fmt.Sprintf("failed generating service: %d with error: %+v", s.idx, s.err)
}
