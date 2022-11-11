package client

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/fx"
	"net"
	"net/http"
	"time"
)

var Module = fx.Module("http/client",
	fx.Provide(
		NewConfig,
		fx.Annotate(
			NewHTTPClient,
			fx.ResultTags(`name:"http/client.base"`),
		),
	),
)

type Config struct {
	*NestedConfig `env:",prefix=HTTP_CLIENT_"`
}

type NestedConfig struct {
	Transport *TransportConfig `env:",prefix=TRANSPORT_"`
}

type TransportConfig struct {
	Dial                  *DialConfig   `env:",prefix=DIAL_"`
	ExpectContinueTimeout time.Duration `env:"EXPECT_CONTINUE_TIMEOUT,default=1s"`
	IdleConnTimeout       time.Duration `env:"IDLE_CONN_TIMEOUT,default=90s"`
	MaxIdleConns          int           `env:"MAX_IDLE_CONNS,default=100"`
	MaxIdleConnsPerHost   int           `env:"MAX_IDLE_CONNS_PER_HOST,default=10"`
	ResponseHeaderTimeout time.Duration `env:"RESPONSE_HEADER_TIMEOUT,default=15s"`
	TLSHandshakeTimeout   time.Duration `env:"TLS_HANDSHAKE_TIMEOUT,default=5s"`
}

type DialConfig struct {
	Timeout   time.Duration `env:"TIMEOUT,default=10s"`
	KeepAlive time.Duration `env:"KEEP_ALIVE,default=30s"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewHTTPClient creates an HTTP client with a custom New Relic transport that
// instruments external requests and adds distributed tracing headers.
func NewHTTPClient(cfg *Config) *http.Client {
	// We do not use DefaultClient since it has no timeout by default:
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779
	httpClient := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: cfg.Transport.ResponseHeaderTimeout,
			Proxy:                 http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				KeepAlive: cfg.Transport.Dial.KeepAlive,
				Timeout:   cfg.Transport.Dial.Timeout,
			}).DialContext,
			MaxIdleConns:          cfg.Transport.MaxIdleConns,
			IdleConnTimeout:       cfg.Transport.IdleConnTimeout,
			TLSHandshakeTimeout:   cfg.Transport.TLSHandshakeTimeout,
			MaxIdleConnsPerHost:   cfg.Transport.MaxIdleConnsPerHost,
			ExpectContinueTimeout: cfg.Transport.ExpectContinueTimeout,
		},
	}
	return httpClient
}
