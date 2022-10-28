package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"net/http"
)

type ModuleOptions struct {
	InvokeFns []any
}

func CreateModule(opts ModuleOptions) fx.Option {
	return fx.Module("handler",
		fx.Provide(
			func() *ModuleOptions {
				return &opts
			},
			NewConfig,
			NewGinEngine,
		),
		fx.Invoke(
			opts.InvokeFns...,
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

func NewGinEngine(
	lc fx.Lifecycle,
	cfg *Config,
	logger zerolog.Logger,
	nrapp *newrelic.Application,
) *gin.Engine {

	// Create Gin router
	r := gin.New()

	// Instrument requests with New Relic telemetry
	r.Use(nrgin.Middleware(nrapp))

	// Middleware to add the Trace ID to the logger
	//r.Use(injectTraceContextLogger(logger, writer, newRelicCfg, zerologCfg))

	//This should work, but its currently not
	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			SkipPaths: []string{"/health"},
		}),
	)
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.Port)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go func() {
				logger.Info().Str("addr", addr).Msg("Serving GraphQL")

				err := r.Run(addr)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal().Err(err).Msg("server failed")
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
