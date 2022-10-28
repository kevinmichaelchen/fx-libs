package logging

import (
	"context"
	"github.com/rs/zerolog"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	nullLogger   = zerolog.Nop()
)

// ExtractNamed is syntactic sugar for extracting a logger from context and
// adding a "logger" field (for named loggers). Similar to what Zap does here:
// https://pkg.go.dev/go.uber.org/zap#Logger.Named
func ExtractNamed(ctx context.Context, name string) *zerolog.Logger {
	l := Extract(ctx)
	ll := l.With().Str("logger", name).Logger()
	return &ll
}

func Extract(ctx context.Context) *zerolog.Logger {
	l, ok := ctx.Value(ctxMarkerKey).(zerolog.Logger)
	//if !ok || l == nil {
	if !ok {
		return &nullLogger
	}
	return &l
}

// ToContext adds the zap.Logger to the context for extraction later.
// Returning the new context that has been created.
func ToContext(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, logger)
}
