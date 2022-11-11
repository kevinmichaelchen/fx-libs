package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"os"
)

var Module = fx.Module("newrelic",
	fx.Provide(
		NewConfig,
		NewNewRelicApplication,
	),
)

type Config struct {
	// e.g., "platform.duo.ngp.ou.manager"
	ServiceName    string `env:"SERVICE_NAME,required"`
	ServiceVersion string `env:"SERVICE_VERSION,default=not-found"`

	// Useful for grouping
	Team string `env:"TEAM"`

	// e.g., "local"
	Env string `env:"ENV,required"`

	*NestedConfig `env:",prefix=NEW_RELIC_"`
}

type NestedConfig struct {
	Enabled bool   `env:"ENABLED,default=true"`
	Key     string `env:"KEY"`

	DebugLogging bool `env:"DEBUG_LOGGING"`

	// Whether the New Relic agent forwards logs to New Relic.
	//
	// The default is false.
	//
	// On all our environments, where we use Kubernetes, we have a dedicated pod
	// that forwards all stdout streams to New Relic, and hence we don't want
	// the in-process New Relic Go Agent to also be forwarding logs.
	//
	// On our local laptops, we may or may not want logs to be forwarded to New
	// Relic. Overall, having the Go Agent perform log forwarding is the
	// exception, not the rule.
	ForwardLogs bool `env:"FORWARD_LOGS,default=false"`

	// Whether logs are decorated for use with the context plugins.
	// Log decoration appends a string to your JSON structured log which is then
	// sent to New Relic for ingestion.
	DecorateLogs bool `env:"DECORATE_LOGS,default=false"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewNewRelicApplication(cfg *Config) (*newrelic.Application, error) {
	enabled := cfg.Enabled
	if enabled && len(cfg.Key) == 0 {
		return nil, errors.New("missing key for New Relic")
	}

	if !enabled {
		fmt.Println("[WARN] New Relic monitoring is not enabled.")
	}

	opts := buildNewRelicConfigOptions(cfg)

	app, err := newrelic.NewApplication(opts...)

	if err != nil {
		return app, err
	}

	return app, nil
}

func buildNewRelicConfigOptions(cfg *Config) []newrelic.ConfigOption {
	enabled := cfg.Enabled
	env := cfg.Env

	labels := map[string]string{
		"team":    cfg.Team,
		"service": cfg.ServiceName,
		"env":     env,
		"version": cfg.ServiceVersion,
	}

	opts := []newrelic.ConfigOption{
		// A base set of configs read from the environment.
		// Latter ConfigOptions may overwrite the Config fields already set.
		newrelic.ConfigFromEnvironment(),
		newrelic.ConfigAppName(cfg.ServiceName + "-" + env),
		newrelic.ConfigLicense(cfg.Key),
		newrelic.ConfigEnabled(enabled),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(
			cfg.ForwardLogs,
		),
		newrelic.ConfigAppLogDecoratingEnabled(cfg.DecorateLogs),
		func(cfg *newrelic.Config) {
			cfg.ErrorCollector.RecordPanics = true
			cfg.CrossApplicationTracer.Enabled = false // this is legacy and is now  DistributedTracerEnabled
			cfg.CustomInsightsEvents.Enabled = true
			cfg.TransactionTracer.Attributes.Enabled = true
			cfg.Labels = labels
		},
	}

	debugLoggerEnabled := enabled && cfg.DebugLogging
	if debugLoggerEnabled {
		opts = append(opts, newrelic.ConfigDebugLogger(os.Stdout))
	}

	return opts
}
