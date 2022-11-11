package newrelic

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
)

func NewRequest(ctx context.Context, r *http.Request) *http.Request {
	txn := newrelic.FromContext(ctx)
	return newrelic.RequestWithTransactionContext(r, txn)
}
