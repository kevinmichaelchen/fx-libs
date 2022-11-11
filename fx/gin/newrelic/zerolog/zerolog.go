package zerolog

import (
	"github.com/gin-gonic/gin"
	pkgNewRelic "github.com/kevinmichaelchen/fx-libs/fx/newrelic"
	pkgZerolog "github.com/kevinmichaelchen/fx-libs/fx/zerolog"
	"github.com/newrelic/go-agent/v3/integrations/logcontext"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Module("gin/newrelic/zerolog",
	fx.Provide(
		fx.Annotate(
			NewGinHandler,
			fx.ResultTags(`name:"gin_newrelic_zerolog_handler"`),
		),
	),
)

func NewGinHandler(
	logger *zerolog.Logger,
	zlw *zerologWriter.ZerologWriter,
	zCfg *pkgZerolog.Config,
	nrCfg *pkgNewRelic.Config,
) gin.HandlerFunc {

	useNewRelic := zCfg.Format == pkgZerolog.FormatJSON &&
		nrCfg.Enabled && nrCfg.ForwardLogs

	return func(c *gin.Context) {
		// Get the New Relic Transaction from the Gin context
		txn := nrgin.Transaction(c)

		var newLogger zerolog.Logger

		// Output duplicates the current logger and sets w as its output.
		// Always create a new logger in order to avoid changing the context of
		// the logger for other threads that may be logging external to this
		// transaction.
		if useNewRelic {
			newLogger = logger.Output(zlw.WithTransaction(txn))
		} else {
			// TODO I THINK this is a safe way of cloning
			// could also try newLogger = logger.Output(os.Stderr)
			newLogger = *logger
		}

		md := txn.GetLinkingMetadata()
		sublogger := newLogger.With().
			Str(logcontext.KeyTraceID, md.TraceID).
			Str(logcontext.KeySpanID, md.SpanID).
			Str(logcontext.KeyEntityName, md.EntityName).
			Str(logcontext.KeyEntityType, md.EntityType).
			Str(logcontext.KeyEntityGUID, md.EntityGUID).
			Str(logcontext.KeyHostname, md.Hostname).
			Logger()

		newCtx := sublogger.WithContext(c)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}
