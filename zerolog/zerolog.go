package zerolog

import (
	"context"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"os"
)

const (
	FormatPretty = "pretty"
	FormatJSON   = "json"
)

var Module = fx.Module("logging",
	fx.Provide(
		NewConfig,
		NewZerolog,
		NewZerologWriter,
	),
)

type Config struct {
	*NestedConfig `env:",prefix=TS_LOG_"`
}

func (c Config) CanForward() bool {
	// only forward logs when it's JSON
	return c.Format == FormatJSON
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

func NewZerologWriter(nrapp *newrelic.Application) zerologWriter.ZerologWriter {
	return zerologWriter.New(os.Stdout, nrapp)
}

func NewZerolog(
	cfg *Config,
	writer zerologWriter.ZerologWriter,
) (zerolog.Logger, error) {
	// Create base logger
	var logger zerolog.Logger
	if cfg.Format == FormatPretty {
		logger = newPrettyLogger(cfg)
	} else {
		logger = zerolog.New(writer)
	}

	logger = logger.With().
		// Add a "time" key (local time as UNIX timestamp)
		Timestamp().
		// Add a "caller" key with the file:line of the caller
		Caller().
		Logger()

	return logger, nil
}
