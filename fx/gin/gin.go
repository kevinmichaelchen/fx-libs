package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"log"
	"net/http"
)

type ModuleOptions struct {
	Invocations []any
}

func CreateModule(opts ModuleOptions) fx.Option {
	return fx.Module("gin",
		fx.Provide(
			NewConfig,
			NewGinEngine,
		),
		fx.Invoke(
			opts.Invocations...,
		),
	)
}

type Config struct {
	Port int `env:"PORT,default=8081"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type NewGinEngineInput struct {
	fx.In

	LC fx.Lifecycle

	Cfg *Config

	// NewRelicHandler is a Gin HandlerFunc (middleware) which instruments the
	// entire Gin handler with New Relic traces.
	NewRelicHandler gin.HandlerFunc `name:"ginnewrelic" optional:"true"`

	// NewRelicZerologHandler is a Gin HandlerFunc (middleware) which injects a
	// trace-aware Zerolog logger into the request context that will collect log
	// metrics, forward logs, and enrich logs depending on how your New Relic
	// application is configured.
	NewRelicZerologHandler gin.HandlerFunc `name:"gin_newrelic_zerolog_handler" optional:"true"`
}

func NewGinEngine(in NewGinEngineInput) *gin.Engine {
	// Create Gin router
	r := gin.New()

	// Instrument requests with New Relic telemetry
	if in.NewRelicHandler != nil {
		r.Use(in.NewRelicHandler)
	}

	// Middleware to add the Trace ID to the logger
	if in.NewRelicZerologHandler != nil {
		r.Use(in.NewRelicZerologHandler)
	}

	// TODO let consumers configure this
	// Don't log the /health endpoint since k8s probes will constantly be
	// hitting it
	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/health"},
		}),
	)

	// Recovery middleware recovers from any panics and writes a 500 if there
	// was one.
	r.Use(gin.Recovery())

	in.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", in.Cfg.Port)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				log.Printf("Serving GraphQL on %s\n", addr)

				err := r.Run(addr)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// TODO see https://github.com/gin-gonic/examples/tree/master/graceful-shutdown/graceful-shutdown
			return nil
		},
	})

	return r
}
