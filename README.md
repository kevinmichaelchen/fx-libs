# fx-libs

[FX](https://github.com/uber-go/fx) modules.

## Modules
As much as possible, we try to keep these modules independent of each other.
Sometimes, when they become glued together (e.g., when New Relic fuses with 
Zerolog), it's convenient to reference a shared configuration struct.

| Module                                       | Description                                     | Dependent on other modules |
|----------------------------------------------|-------------------------------------------------|----------------------------|
| [`client`](./client)                         | New Relic-instrumented HTTP client              | NO                         |
| [`genqlient`](./genqlient)                   | GraphQL client built on genqlient               | NO                         |
| [`gin`](./gin)                               | Gin HTTP server                                 | NO                         |
| [`ginnewrelic`](./ginnewrelic)               | Gin middleware for New Relic instrumentation    | NO                         |
| [`ginnewreliczerolog`](./ginnewreliczerolog) | Gin middleware for New Relic "Logs in Context"  | YES                        |
| [`gqlgen`](./gqlgen)                         | Interceptors for gqlgen                         | NO                         |
| [`newrelic`](./newrelic)                     | Configures a New Relic Go Agent                 | NO                         |
| [`zerolog`](./zerolog)                       | Configures a New Relic Zerolog logger           | NO                         |

## Installation
```shell
go get -v go.uber.org/fx

# Install only the modules you wish to use
go get -v \
  bitbucket.org/tskevinchen/fx-libs/client \
  bitbucket.org/tskevinchen/fx-libs/genqlient \
  bitbucket.org/tskevinchen/fx-libs/gin \
  bitbucket.org/tskevinchen/fx-libs/ginnewrelic \
  bitbucket.org/tskevinchen/fx-libs/ginnewreliczerolog \
  bitbucket.org/tskevinchen/fx-libs/gqlgen \
  bitbucket.org/tskevinchen/fx-libs/newrelic \
  bitbucket.org/tskevinchen/fx-libs/zerolog
```

### Updating dependencies
```shell
go list all | grep bitbucket.org/tskevinchen/fx-libs | xargs go get -v
```

## Usage
See the [example project](./example).

In `main.go`:
```go
package main

import (
	"bitbucket.org/tskevinchen/fx-libs/client"
	"bitbucket.org/tskevinchen/fx-libs/genqlient"
	"bitbucket.org/tskevinchen/fx-libs/gin"
	"bitbucket.org/tskevinchen/fx-libs/ginnewrelic"
	"bitbucket.org/tskevinchen/fx-libs/ginnewreliczerolog"
	"bitbucket.org/tskevinchen/fx-libs/newrelic"
	"bitbucket.org/tskevinchen/fx-libs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Options(
	// New Relic-instrumented HTTP client
	client.Module,

	// genqlient GraphQL handler
	genqlient.Module,

	// Gin HTTP handler
	gin.CreateModule(gin.ModuleOptions{
		Invocations: []any{
			// Add a registration function here to register your business logic
			// layer (often called the "Service layer" or "Use case layer") to
			// your GraphQL Handler.
		},
	}),

	// Gin middleware for New Relic instrumentation
	ginnewrelic.Module,

	// Gin middleware for trace-aware log forwarding
	ginnewreliczerolog.Module,

	// New Relic Go Agent
	newrelic.Module,

	// Zerolog logging
	zerolog.Module,
)

func main() {
	a := fx.New(
		Module,
	)
	a.Run()
}
```

## Contributing

### Tidying
```shell
go work sync
```
