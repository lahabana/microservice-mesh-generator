# service-mesh-generator

[![Run in Insomnia}](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=microservice-mesh-generator&uri=https%3A%2F%2Fraw.githubusercontent.com%2Flahabana%2Fmicroservice-mesh-generator%2Fmain%2Fopenapi.yaml)

A tool to easily build meshes of services. It is hosted [there](https://mservice-mesh-generator-lahabana.koyeb.app/).

When working on [Kuma](https://kuma.io) it is used to write e2e tests and perf tests, demo, experiment... 

## Features

- Generate a random mesh of services which talk to each other.
- Generate a kubernetes manifest to run this mesh either with [fake-service](https://github.com/nicholasjackson/fake-service) or [api-play](https://github.com/lahabana/api-play).
- Reproducible setups 

## Usage

### Hosted version

We host a version using [Koyeb](https://koyeb.com) it's available [there](https://mservice-mesh-generator-lahabana.koyeb.app/).

### Command line

```shell
docker run --rm ghcr.io/lahabana/microservice-mesh-generator:main --help
```

You can for example create a random mesh with:

```shell
docker run --rm -p 8080:8080 ghcr.io/lahabana/microservice-mesh-generator:main -output k8s | kubectl apply -f -
```

and then access the first service with:

```shell
kubectl port-forward -n microservice-mesh svc/api-play-000 8080:8080
```

You can then access: http://localhost:8080/api/dynamic/microservice_mesh

### Local server

```shell
docker run --rm -p 8080:8080 ghcr.io/lahabana/microservice-mesh-generator:main -server
```

Access at `http://localhost:8080`

### As a library

The code to generate things can be used as a library.
With this library you can:

- define [a mesh](https://github.com/lahabana/microservice-mesh-generator/blob/a1290d4e7c39cad26dac113fd74758578179ab73/pkg/generators/k8s/generic_test.go#L16-L23)
- use the [random generator](https://github.com/lahabana/microservice-mesh-generator/blob/a1290d4e7c39cad26dac113fd74758578179ab73/main.go#L36)
- just define your own [app generator](https://github.com/lahabana/microservice-mesh-generator/blob/main/pkg/generators/k8s/apiplay/generator.go)

## TODO

- add a way to define your own mesh
- add a library of realistic meshes
- add a checkbox to add kuma params
- add a checkbox to add istio params
- add a checkbox to add linkerd params
- have a better domain
- Options to add latency/errors etc.
