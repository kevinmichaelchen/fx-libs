package newrelic

import (
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
	"io"
	"os"
)

var Module = fx.Module("zerolog/newrelic",
	fx.Provide(
		// We're using the NewZerologWriter provider twice.
		// The first time, we provide it as *zerologWriter.ZerologWriter.
		NewZerologWriter,
		fx.Annotate(
			// The second time, we provide it as an io.Writer.
			// This allows the ./zerolog module to decouple entirely from
			// anything related to New Relic.
			NewZerologWriter,
			fx.As(new(io.Writer)),
			fx.ResultTags(`name:"zerolog_newrelic_writer"`),
		),
	),
)

func NewZerologWriter(nrapp *newrelic.Application) *zerologWriter.ZerologWriter {
	zlWriter := zerologWriter.New(os.Stdout, nrapp)
	return &zlWriter
}
