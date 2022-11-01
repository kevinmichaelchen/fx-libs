package gqlgen

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewRelicMiddleware() graphql.OperationMiddleware {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		// e.g., "mutation-CreateFoo"
		name := fmt.Sprintf("%s-%s", oc.Operation.Operation, oc.OperationName)
		txn := newrelic.FromContext(ctx)
		txn.SetName(name)
		return next(ctx)
	}
}
