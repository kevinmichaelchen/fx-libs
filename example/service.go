package main

import (
	"context"
	"github.com/kevinmichaelchen/fx-libs/fx/http/client/newrelic"
	zl "github.com/rs/zerolog"
	"io"
	"net/http"
)

type Service struct {
	httpClient *http.Client
}

func (s *Service) OK(ctx context.Context) {
	logger := zl.Ctx(ctx)
	logger.Info().Dur("my.http.client.timeout", s.httpClient.Timeout).Msg("http client")
	logger.Info().Msg("inside /ok")

	req, err := http.NewRequest(http.MethodGet, "http://worldtimeapi.org/api/timezone/America/New_York", nil)
	if err != nil {
		logger.Err(err).Msg("failed to build http request")
		return
	}

	// Add the transaction to the request context. This step is required for B3
	// header propagation: https://github.com/newrelic/go-agent/blob/master/v3/integrations/nrb3/example_test.go
	req = newrelic.NewRequest(ctx, req)

	res, err := s.httpClient.Do(req)
	if err != nil {
		logger.Err(err).Msg("request failed")
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Err(err).Msg("failed to read response body")
		return
	}

	logger.Info().RawJSON("res", b).Msg("res")
}
