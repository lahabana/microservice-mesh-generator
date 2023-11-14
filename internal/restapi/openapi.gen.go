// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package restapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Details           string              `json:"details"`
	InvalidParameters *[]InvalidParameter `json:"invalid_parameters,omitempty"`
	Status            int                 `json:"status"`
}

// Health defines model for Health.
type Health struct {
	Status int `json:"status"`
}

// HomeResponse defines model for HomeResponse.
type HomeResponse struct {
	Commit   string `json:"commit"`
	Hostname string `json:"hostname"`
	Target   string `json:"target"`
	Version  string `json:"version"`
}

// InvalidParameter defines model for InvalidParameter.
type InvalidParameter struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

// RandomMeshResponse defines model for RandomMeshResponse.
type RandomMeshResponse = map[string]interface{}

// GenerateRandomParams defines parameters for GenerateRandom.
type GenerateRandomParams struct {
	// K8sApp The app to use
	K8sApp *string `form:"k8sApp,omitempty" json:"k8sApp,omitempty"`

	// K8sNamespace the namespace to use
	K8sNamespace *string `form:"k8sNamespace,omitempty" json:"k8sNamespace,omitempty"`

	// Seed the seed to use for deterministic randomness
	Seed *int64 `form:"seed,omitempty" json:"seed,omitempty"`

	// K8s whether or not to return kubernetes manifest
	K8s *bool `form:"k8s,omitempty" json:"k8s,omitempty"`

	// NumServices number of services to run
	NumServices *int `form:"numServices,omitempty" json:"numServices,omitempty"`

	// MinReplicas minimum number of replicas per service
	MinReplicas *int `form:"minReplicas,omitempty" json:"minReplicas,omitempty"`

	// MaxReplicas maximum number of replicas per service
	MaxReplicas *int `form:"maxReplicas,omitempty" json:"maxReplicas,omitempty"`

	// PercentEdge maximum number of replicas per service
	PercentEdge *int `form:"percentEdge,omitempty" json:"percentEdge,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// home
	// (GET /api)
	Home(c *gin.Context)
	// generate a random mesh
	// (GET /api/random.{format})
	GenerateRandom(c *gin.Context, format string, params GenerateRandomParams)
	// healthcheck
	// (GET /health)
	Health(c *gin.Context)
	// healthcheck
	// (GET /ready)
	Ready(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// Home operation middleware
func (siw *ServerInterfaceWrapper) Home(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Home(c)
}

// GenerateRandom operation middleware
func (siw *ServerInterfaceWrapper) GenerateRandom(c *gin.Context) {

	var err error

	// ------------- Path parameter "format" -------------
	var format string

	err = runtime.BindStyledParameter("simple", false, "format", c.Param("format"), &format)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter format: %w", err), http.StatusBadRequest)
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GenerateRandomParams

	// ------------- Optional query parameter "k8sApp" -------------

	err = runtime.BindQueryParameter("form", true, false, "k8sApp", c.Request.URL.Query(), &params.K8sApp)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter k8sApp: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "k8sNamespace" -------------

	err = runtime.BindQueryParameter("form", true, false, "k8sNamespace", c.Request.URL.Query(), &params.K8sNamespace)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter k8sNamespace: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "seed" -------------

	err = runtime.BindQueryParameter("form", true, false, "seed", c.Request.URL.Query(), &params.Seed)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter seed: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "k8s" -------------

	err = runtime.BindQueryParameter("form", true, false, "k8s", c.Request.URL.Query(), &params.K8s)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter k8s: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "numServices" -------------

	err = runtime.BindQueryParameter("form", true, false, "numServices", c.Request.URL.Query(), &params.NumServices)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter numServices: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "minReplicas" -------------

	err = runtime.BindQueryParameter("form", true, false, "minReplicas", c.Request.URL.Query(), &params.MinReplicas)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter minReplicas: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "maxReplicas" -------------

	err = runtime.BindQueryParameter("form", true, false, "maxReplicas", c.Request.URL.Query(), &params.MaxReplicas)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter maxReplicas: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "percentEdge" -------------

	err = runtime.BindQueryParameter("form", true, false, "percentEdge", c.Request.URL.Query(), &params.PercentEdge)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter percentEdge: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GenerateRandom(c, format, params)
}

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Health(c)
}

// Ready operation middleware
func (siw *ServerInterfaceWrapper) Ready(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.Ready(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/api", wrapper.Home)
	router.GET(options.BaseURL+"/api/random.:format", wrapper.GenerateRandom)
	router.GET(options.BaseURL+"/health", wrapper.Health)
	router.GET(options.BaseURL+"/ready", wrapper.Ready)
}
