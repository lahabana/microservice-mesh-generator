package random

import (
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"math/rand"
)

// GenerateRandomMesh creates a mesh of some instances with some replicas.
func GenerateRandomMesh(seed int64, numServices, percentEdges, minReplicas, maxReplicas int) apis.ServiceGraph {
	r := rand.New(rand.NewSource(seed))
	srvs := apis.ServiceGraph{}
	for i := 0; i < numServices; i++ {
		numInstances := 1
		if maxReplicas >= minReplicas {
			numInstances = (r.Int() % (1 + maxReplicas - minReplicas)) + minReplicas
		}
		srvs.Services = append(srvs.Services, apis.Service{Idx: i, Replicas: numInstances})
	}
	// That's the whole story of DAG and topological sort with triangular matrix.
	for i := 0; i < numServices; i++ {
		for j := i + 1; j < numServices; j++ {
			if r.Int()%(j-i) == 0 && r.Int()%100 < percentEdges {
				srvs.Services[i].Edges = append(srvs.Services[i].Edges, j)
			}
		}
	}
	return srvs
}
