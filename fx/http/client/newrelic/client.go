package newrelic

import (
	"github.com/newrelic/go-agent/v3/integrations/nrb3"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Module("http/client/newrelic",
	fx.Provide(
		NewHTTPClient,
	),
)

type NewHTTPClientInput struct {
	fx.In

	HTTPClient *http.Client `name:"http/client.base"`
}

// NewHTTPClient creates an HTTP client with a custom New Relic transport that
// instruments external requests and adds distributed tracing headers.
func NewHTTPClient(in NewHTTPClientInput) *http.Client {
	baseClient := in.HTTPClient
	return &http.Client{
		// When defining the client, set the Transport to the NewRoundTripper. This
		// will create ExternalSegments and add B3 headers for each request.
		Transport:     nrb3.NewRoundTripper(baseClient.Transport),
		CheckRedirect: baseClient.CheckRedirect,
		Jar:           baseClient.Jar,
		Timeout:       baseClient.Timeout,
	}
}
