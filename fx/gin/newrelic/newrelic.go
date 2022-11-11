package newrelic

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
)

var Module = fx.Module("gin/newrelic",
	fx.Provide(
		fx.Annotate(
			NewGinHandler,
			fx.ResultTags(`name:"ginnewrelic"`),
		),
	),
)

func NewGinHandler(nrapp *newrelic.Application) gin.HandlerFunc {
	return nrgin.Middleware(nrapp)
}
