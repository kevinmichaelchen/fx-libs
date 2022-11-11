package main

import (
	"errors"
	"fmt"
	gg "github.com/gin-gonic/gin"
	"github.com/ipfans/fxlogger"
	"github.com/kevinmichaelchen/fx-libs/fx/gin"
	"github.com/kevinmichaelchen/fx-libs/fx/recipes/standard"
	"github.com/newrelic/go-agent/v3/newrelic"
	zl "github.com/rs/zerolog"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Options(
	standard.Module,

	fx.Provide(
		func(httpClient *http.Client) *Service {
			return &Service{httpClient: httpClient}
		},
	),

	// Gin HTTP handler
	gin.CreateModule(gin.ModuleOptions{
		Invocations: []any{
			// Add a registration function here to register your business logic
			// layer (often called the "Service layer" or "Use case layer") to
			// your GraphQL Handler.
			RegisterHandler,
		},
	}),
)

func main() {
	a := fx.New(
		Module,
		fx.WithLogger(
			fxlogger.Default(),
		),
	)
	a.Run()
}

func RegisterHandler(svc *Service, r *gg.Engine) {
	// for k8s health check
	r.GET("/health", func(c *gg.Context) {
		c.Writer.Write([]byte("ok"))
	})

	r.GET("/ok", func(c *gg.Context) {
		ctx := c.Request.Context()
		svc.OK(ctx)
		c.Writer.Write([]byte("ok"))
	})

	r.GET("/err", func(c *gg.Context) {
		ctx := c.Request.Context()
		logger := zl.Ctx(ctx)
		txn := newrelic.FromContext(ctx)
		err := errors.New("new internal error")
		logger.Error().Err(err).Msg("uh oh")
		txn.NoticeError(fmt.Errorf("a little bit of context: %w", err))
	})
}
