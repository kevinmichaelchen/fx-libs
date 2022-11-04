package client

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("client",
	fx.Provide(
		NewHTTPClient,
	),
)

// NewHTTPClient creates an HTTP client with a custom New Relic transport that
// instruments external requests and adds distributed tracing headers.
func NewHTTPClient() *http.Client {
	httpClient := http.DefaultClient
	httpClient.Transport = newrelic.NewRoundTripper(httpClient.Transport)
	return httpClient
}
