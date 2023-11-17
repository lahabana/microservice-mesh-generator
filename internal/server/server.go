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

func (s *srv) GenerateRandom(c *gin.Context, format restapi.GenerateRandomParamsFormat, params restapi.GenerateRandomParams) {
	var invParams []restapi.InvalidParameter
	ctx := c.Request.Context()
	config := generate.DefaultConfig()
	if params.K8sApp != nil {
		config.K8sApp = string(*params.K8sApp)
	}
	if params.K8sNamespace != nil {
		config.K8sNamespace = *params.K8sNamespace
	}
	if params.Seed != nil {
		config.Seed = int64(*params.Seed)
	}
	contentType := ""
	switch format {
	case restapi.Empty, restapi.Yaml:
		contentType = "application/yaml"
		if params.K8s != nil && *params.K8s {
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
		percentEdge = 50
	}
	if percentEdge < 0 || percentEdge > 100 {
		invParams = append(invParams, restapi.InvalidParameter{
			Field:  "percentEdge",
			Reason: "must be between 0 and 99",
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
			s.l.InfoContext(ctx, "failed request", "error", err)
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
	restapi.RegisterHandlersWithOptions(engine, &srv{}, restapi.GinServerOptions{
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
