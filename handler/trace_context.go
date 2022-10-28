package handler

//import (
//	newrelic2 "github.com/kevinmichaelchen/fx-libs/newrelic"
//	zerolog2 "github.com/kevinmichaelchen/fx-libs/zerolog"
//	"github.com/newrelic/go-agent/v3/integrations/logcontext"
//	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
//	"github.com/rs/zerolog"
//)
//
//func injectTraceContextLogger(
//	logger zerolog.Logger,
//	writer zerologWriter.ZerologWriter,
//	newRelicCfg *newrelic2.Config,
//	zerologCfg *zerolog2.Config,
//) gin.HandlerFunc {
//
//	return func(c *gin.Context) {
//		// Get the New Relic Transaction from the Gin context
//		txn := nrgin.Transaction(c)
//
//		var newLogger zerolog.Logger
//
//		// Output duplicates the current logger and sets w as its output.
//		// Always create a new logger in order to avoid changing the context of
//		// the logger for other threads that may be logging external to this
//		// transaction.
//		if zerologCfg.CanForward() && newRelicCfg.ShouldForwardLogs() {
//			newLogger = logger.Output(writer.WithTransaction(txn))
//		} else {
//			// TODO I THINK this is a safe way of cloning
//			newLogger = *(&logger)
//		}
//
//		md := txn.GetLinkingMetadata()
//		sublogger := newLogger.With().
//			Str(logcontext.KeyTraceID, md.TraceID).
//			Str(logcontext.KeySpanID, md.SpanID).
//			Str(logcontext.KeyEntityName, md.EntityName).
//			Str(logcontext.KeyEntityType, md.EntityType).
//			Str(logcontext.KeyEntityGUID, md.EntityGUID).
//			Str(logcontext.KeyHostname, md.Hostname).
//			Logger()
//
//		newCtx := logging.ToContext(c, sublogger)
//		c.Request = c.Request.WithContext(newCtx)
//		c.Next()
//	}
//}
