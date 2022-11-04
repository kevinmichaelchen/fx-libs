package main

import (
	gg "github.com/gin-gonic/gin"
	"github.com/kevinmichaelchen/fx-libs/client"
	"github.com/kevinmichaelchen/fx-libs/genqlient"
	"github.com/kevinmichaelchen/fx-libs/gin"
	"github.com/kevinmichaelchen/fx-libs/ginnewrelic"
	"github.com/kevinmichaelchen/fx-libs/ginnewreliczerolog"
	"github.com/kevinmichaelchen/fx-libs/newrelic"
	"github.com/kevinmichaelchen/fx-libs/zerolog"
	zl "github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Options(
	// New Relic-instrumented HTTP client
	client.Module,

	// genqlient GraphQL handler
	genqlient.Module,

	// Gin HTTP handler
	gin.CreateModule(gin.ModuleOptions{
		Invocations: []any{
			// Add a registration function here to register your business logic
			// layer (often called the "Service layer" or "Use case layer") to
			// your GraphQL Handler.
			RegisterHandler,
		},
	}),

	// Gin middleware for New Relic instrumentation
	ginnewrelic.Module,

	// Gin middleware for trace-aware log forwarding
	ginnewreliczerolog.Module,

	// New Relic Go Agent
	newrelic.Module,

	// Zerolog logging
	zerolog.Module,
)

func main() {
	a := fx.New(
		Module,
	)
	a.Run()
}

func RegisterHandler(r *gg.Engine) {
	// for k8s health check
	r.GET("/health", func(c *gg.Context) {
		c.Writer.Write([]byte("ok"))
	})

	r.GET("/ok", func(c *gg.Context) {
		logger := zl.Ctx(c.Request.Context())
		logger.Info().Msg("inside /ok")
		c.Writer.Write([]byte("ok"))
	})
}
