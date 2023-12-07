package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lahabana/microservice-mesh-generator/internal/generate"
	"github.com/lahabana/microservice-mesh-generator/internal/restapi"
	"github.com/lahabana/microservice-mesh-generator/internal/server/www"
	"github.com/lahabana/microservice-mesh-generator/pkg/apis"
	"github.com/lahabana/microservice-mesh-generator/pkg/version"
	"github.com/lahabana/otel-gin/pkg/observability"
	"log/slog"
	"net/http"
	"os"
	"runtime"
)

type srv struct {
	l *slog.Logger
}

func (s *srv) PostApiDefineFormat(c *gin.Context, format restapi.OutputFormat, params restapi.PostApiDefineFormatParams) {
	ctx := c.Request.Context()
	if c.Request.ContentLength > (1 << 20) {
		c.PureJSON(http.StatusBadRequest, restapi.ErrorResponse{
			Status:  http.StatusBadRequest,
			Details: "Bad Request",
			InvalidParameters: &[]restapi.InvalidParameter{
				{
					Field:  "payload",
					Reason: "Max payload size is 1MiB",
				},
			},
		})
		return
	}

	var invParams []restapi.InvalidParameter
	inputGraph := restapi.PostApiDefineFormatJSONRequestBody{}
	if err := c.BindJSON(&inputGraph); err != nil {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "payload",
			Reason: "failed to parse payload: " + err.Error(),
		})
	}

	graph := apis.ServiceGraph{
		GenerationParams: "provided by api",
	}
	if len(inputGraph.Services) == 0 || len(inputGraph.Services) > 5000 {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "payload.services",
			Reason: "can't have 0 or more than 5000 services",
		})
	} else {
		for i, srv := range inputGraph.Services {
			if len(srv.Edges) > 50 {
				invParams = append(invParams, restapi.InvalidParameter{
					Field:  fmt.Sprintf("payload.services[%d].edges", i),
					Reason: fmt.Sprintf("can't have more than 50 edges"),
				})
			} else {
				for j, edge := range srv.Edges {
					if edge > len(inputGraph.Services) || edge < 0 {
						invParams = append(invParams, restapi.InvalidParameter{
							Field:  fmt.Sprintf("payload.services[%d].edges[%d]", i, j),
							Reason: fmt.Sprintf("the destination of the call is not an existing entity, max index %d", len(inputGraph.Services)),
						})
					}
				}
				graph.Services = append(graph.Services, apis.Service{
					Idx:      i,
					Replicas: srv.Replicas,
					Edges:    srv.Edges,
				})
			}
		}
	}

	config, contentType, invConfParams := s.extractConfig(format, params.K8s, params.K8sApp, params.K8sNamespace, nil)
	invParams = append(invParams, invConfParams...)

	if len(invParams) > 0 {
		c.PureJSON(http.StatusBadRequest, restapi.ErrorResponse{
			Status:            http.StatusBadRequest,
			Details:           "Bad Request",
			InvalidParameters: &invParams,
		})
		return
	}

	buf := bytes.Buffer{}
	config.Writer = &buf
	err := generate.Run(config, func(seed int64) (apis.ServiceGraph, error) {
		return graph, nil
	})
	if err != nil {
		if errors.Is(err, &generate.InvalidConfError{}) {
			c.PureJSON(http.StatusBadRequest, restapi.ErrorResponse{
				Status:  http.StatusBadRequest,
				Details: err.Error(),
			})
			return
		} else {
			s.l.ErrorContext(ctx, "failed request", "error", err)
			c.PureJSON(http.StatusInternalServerError, restapi.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Details: "internal error",
			})
			return
		}
	}
	c.Data(http.StatusOK, contentType, buf.Bytes())

}

func (s *srv) extractConfig(format restapi.OutputFormat, k8s *bool, k8sApp *restapi.K8sAppType, k8sNamespace *string, seed *int) (generate.Config, string, []restapi.InvalidParameter) {
	s.l.Info("foo", "format", format)
	var invParams []restapi.InvalidParameter
	config := generate.DefaultConfig()
	if k8sApp != nil {
		config.K8sApp = string(*k8sApp)
	}
	if k8sNamespace != nil {
		config.K8sNamespace = *k8sNamespace
	}
	contentType := ""
	switch format {
	case restapi.Empty, restapi.Yaml:
		contentType = "application/yaml"
		if k8s != nil && *k8s {
			config.Output = "k8s"
		} else {
			config.Output = "yaml"
		}
	case restapi.Mmd:
		contentType = "text/vnd.mermaid"
		config.Output = "mermaid"
	case restapi.Gv:
		contentType = "text/vnd.graphviz"
		config.Output = "dot"
	case restapi.Json:
		contentType = "application/json"
		config.Output = "json"
	default:
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "format",
			Reason: "not a supported type",
		})
	}
	if seed != nil {
		config.Seed = int64(*seed)
	}
	return config, contentType, invParams
}

func (s *srv) GenerateRandom(c *gin.Context, format restapi.OutputFormat, params restapi.GenerateRandomParams) {
	var invParams []restapi.InvalidParameter
	ctx := c.Request.Context()
	config, contentType, invConfParams := s.extractConfig(format, params.K8s, params.K8sApp, params.K8sNamespace, params.Seed)
	invParams = append(invParams, invConfParams...)
	numServices := 5
	if params.NumServices != nil {
		numServices = *params.NumServices
	}
	if numServices <= 0 {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "numServices",
			Reason: "can't be null or negative",
		})
	}
	percentEdge := 50
	if params.PercentEdge != nil {
		percentEdge = *params.PercentEdge
	}
	if percentEdge < 0 || percentEdge > 100 {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "percentEdge",
			Reason: "must be between 0 and 100",
		})
	}
	minReplicas := 2
	if params.MinReplicas != nil {
		minReplicas = *params.MinReplicas
	}
	if minReplicas <= 0 {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "minReplicas",
			Reason: "must > 0",
		})
	}
	maxReplicas := minReplicas
	if params.MaxReplicas != nil {
		maxReplicas = *params.MaxReplicas
	}
	if maxReplicas < minReplicas {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "maxReplicas",
			Reason: "can't be lower than minReplicas",
		})
	}
	if len(invParams) > 0 {
		c.PureJSON(http.StatusBadRequest, restapi.ErrorResponse{
			Status:            http.StatusBadRequest,
			Details:           "Bad Request",
			InvalidParameters: &invParams,
		})
		return
	}

	buf := bytes.Buffer{}
	config.Writer = &buf
	err := generate.Run(config, func(seed int64) (apis.ServiceGraph, error) {
		graph := apis.GenerateRandomMesh(seed, numServices, percentEdge, minReplicas, maxReplicas)
		return graph, nil
	})
	if err != nil {
		if errors.Is(err, &generate.InvalidConfError{}) {
			c.PureJSON(http.StatusBadRequest, restapi.ErrorResponse{
				Status:  http.StatusBadRequest,
				Details: err.Error(),
			})
			return
		} else {
			s.l.ErrorContext(ctx, "failed request", "error", err)
			c.PureJSON(http.StatusInternalServerError, restapi.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Details: "internal error",
			})
			return
		}
	}
	c.Data(http.StatusOK, contentType, buf.Bytes())
}

func (s *srv) Home(c *gin.Context) {
	host, _ := os.Hostname()
	c.PureJSON(http.StatusOK, restapi.HomeResponse{
		Hostname: host,
		Version:  version.Version,
		Commit:   version.Commit,
		Target:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})
}

func (s *srv) Health(c *gin.Context) {
	c.PureJSON(http.StatusOK, restapi.Health{Status: http.StatusOK})
}

func (s *srv) Ready(c *gin.Context) {
	c.PureJSON(http.StatusOK, restapi.Health{Status: http.StatusOK})
}

func Start(ctx context.Context) error {
	obs, err := observability.Init(ctx, "api-play", slog.LevelDebug, observability.OTLPNone, observability.OTLPNone)
	if err != nil {
		panic(err)
	}
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		obs.Logger().ErrorContext(c.Request.Context(), "panic in handler", "error", err)

		c.PureJSON(http.StatusInternalServerError, &restapi.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Details: "Internal server error",
		})
	}), obs.Middleware())
	restapi.RegisterHandlersWithOptions(engine, &srv{l: obs.Logger().WithGroup("handler")}, restapi.GinServerOptions{
		ErrorHandler: func(c *gin.Context, err error, status int) {
			c.PureJSON(status, &restapi.ErrorResponse{
				Status:  status,
				Details: err.Error(),
			})
		},
	})
	hfs := http.FS(www.Content)
	engine.StaticFileFS("/", "index.htm", hfs)
	engine.StaticFS("/static", hfs)
	return engine.Run()
}
