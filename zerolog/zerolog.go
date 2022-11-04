package zerolog

import (
	"context"
	"errors"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"os"
	"time"
)

const (
	FormatPretty = "pretty"
	FormatJSON   = "json"
)

var Module = fx.Module("zerolog",
	fx.Provide(
		NewConfig,
		NewZerolog,
		NewZerologWriter,
	),
)

type Config struct {
	*NestedConfig `env:",prefix=TS_LOG_"`
}

type NestedConfig struct {
	Level  string `env:"LEVEL,required"`
	Format string `env:"FORMAT,required"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewZerolog(cfg *Config, zlw *zerologWriter.ZerologWriter) (*zerolog.Logger, error) {
	// Create base logger
	var logger zerolog.Logger
	if cfg.Format == FormatPretty {
		logger = log.Output(&zerolog.ConsoleWriter{
			Out: os.Stderr,
			// We always want colors when running in pretty mode.
			NoColor: false,
			// default time format is time.Kitchen (e.g., "3:04PM")
			TimeFormat: time.RFC3339Nano,
		})
	} else if cfg.Format == FormatJSON {
		logger = zerolog.New(zlw)
	} else {
		return nil, errors.New("unknown format")
	}

	logger = logger.With().
		// Add a "time" key (local time as UNIX timestamp)
		Timestamp().
		// Add a "caller" key with the file:line of the caller
		Caller().
		Logger()

	return &logger, nil
}

func NewZerologWriter(nrapp *newrelic.Application) *zerologWriter.ZerologWriter {
	zlWriter := zerologWriter.New(os.Stdout, nrapp)
	return &zlWriter
}
