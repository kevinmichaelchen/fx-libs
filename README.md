# fx-libs

[FX](https://github.com/uber-go/fx) modules.

## Installation
```shell
go get -v go.uber.org/fx

# Install only the modules you wish to use
go get -v \
  github.com/kevinmichaelchen/fx-libs/client \
  github.com/kevinmichaelchen/fx-libs/gqlgen \
  github.com/kevinmichaelchen/fx-libs/handler \
  github.com/kevinmichaelchen/fx-libs/newrelic \
  github.com/kevinmichaelchen/fx-libs/zerolog
```

### Updating dependencies
```shell
go list all | grep github.com/kevinmichaelchen/fx-libs | xargs go get -v
```

## Usage
In `main.go`:
```go
package main

import (
	"github.com/kevinmichaelchen/fx-libs/client"
	"github.com/kevinmichaelchen/fx-libs/handler"
	"github.com/kevinmichaelchen/fx-libs/newrelic"
	"github.com/kevinmichaelchen/fx-libs/zerolog"
	"go.uber.org/fx"
)

var Module = fx.Options(
	client.Module,
	handler.CreateModule(handler.ModuleOptions{
		Invocations: []any{
			// Add a registration function here to register your business logic
			// layer (often called the "Service layer" or "Use case layer" to
			// your GraphQL Handler.
		},
		Providers: []any{
			func(zCfg *zerolog.Config, nrCfg *newrelic.Config) *handler.UseNewRelicOutput {
				return &handler.UseNewRelicOutput{
					UseNewRelic: zCfg.Format == zerolog.FormatJSON && nrCfg.Enabled && nrCfg.ForwardLogs,
				}
			},
		},
	}),
	newrelic.Module,
	zerolog.Module,
)

func main() {
	a := fx.New(
		Module,
	)
	a.Run()
}
```

