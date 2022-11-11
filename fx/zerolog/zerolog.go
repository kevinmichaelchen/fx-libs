package zerolog

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"io"
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
		fx.Annotate(
			NewZerolog,
			fx.ParamTags(
				// the cfg param isn't annotated
				``,
				// the io.Writer param is annotated
				`name:"zerolog_newrelic_writer" optional:"true"`),
		),
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

func NewZerolog(cfg *Config, nrw io.Writer) (*zerolog.Logger, error) {
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
		logger = zerolog.New(nrw)
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
