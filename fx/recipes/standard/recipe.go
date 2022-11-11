package standard

import (
	"github.com/kevinmichaelchen/fx-libs/fx/genqlient"
	pkgGinNR "github.com/kevinmichaelchen/fx-libs/fx/gin/newrelic"
	pkgGinNRZL "github.com/kevinmichaelchen/fx-libs/fx/gin/newrelic/zerolog"
	pkgHTTPClient "github.com/kevinmichaelchen/fx-libs/fx/http/client"
	pkgNewRelicHTTPClient "github.com/kevinmichaelchen/fx-libs/fx/http/client/newrelic"
	"github.com/kevinmichaelchen/fx-libs/fx/newrelic"
	"github.com/kevinmichaelchen/fx-libs/fx/zerolog"
	pkgZLNR "github.com/kevinmichaelchen/fx-libs/fx/zerolog/newrelic"
	"go.uber.org/fx"
)

var Module = fx.Module("recipe_standard",
	// New Relic-instrumented HTTP client
	pkgNewRelicHTTPClient.Module,
	pkgHTTPClient.Module,

	// genqlient GraphQL handler
	genqlient.Module,

	// Gin middleware for New Relic instrumentation
	pkgGinNR.Module,

	// Gin middleware for trace-aware log forwarding
	pkgGinNRZL.Module,

	// New Relic Go Agent
	newrelic.Module,

	// Zerolog logging
	zerolog.Module,
	pkgZLNR.Module,
)
