package zerolog

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func newPrettyLogger(cfg *Config) zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
		// We always want colors when running in pretty mode.
		NoColor: false,
		// default time format is time.Kitchen (e.g., "3:04PM")
		TimeFormat: time.RFC3339Nano,
	})
}
