# fx-libs

This project is a [Multi-Module Go 1.18 Workspace](https://go.dev/doc/tutorial/workspaces)
that contains mostly [FX](https://github.com/uber-go/fx) modules for our team's
microservices.

If you've never used FX, we highly recommend it. [Preslav Mihaylov](https://youtu.be/i2DN-f6rPL8?t=330)
does a brilliant job in explaining what it is, why you'd want to use it, and the
trade-offs teams should consider before adopting it.

## Modules
As much as possible, we try to keep these modules independent of each other.
Sometimes, when they become glued together (as is the case in
[`gin/newrelic/zerolog`](fx/gin/newrelic/zerolog)), it's convenient to reference
a shared configuration struct.

| Module                                               | Description                                    | Dependent on other modules |
|------------------------------------------------------|------------------------------------------------|----------------------------|
| [`fx/genqlient`](fx/genqlient)                       | GraphQL client built on genqlient              | NO                         |
| [`fx/gin`](fx/gin)                                   | Gin HTTP server                                | NO                         |
| [`fx/gin/newrelic`](fx/gin/newrelic)                 | Gin middleware for New Relic instrumentation   | NO                         |
| [`fx/gin/newrelic/zerolog`](fx/gin/newrelic/zerolog) | Gin middleware for New Relic "Logs in Context" | YES                        |
| [`fx/http/client`](fx/http/client)                   | Fine-tuned HTTP client                         | NO                         |
| [`fx/http/client/newrelic`](fx/http/client/newrelic) | New Relic-instrumented HTTP client             | NO                         |
| [`fx/newrelic`](fx/newrelic)                         | Configures a New Relic Go Agent                | NO                         |
| [`fx/zerolog`](fx/zerolog)                           | Configures a New Relic Zerolog logger          | YES                        |
| [`gqlgen`](./gqlgen)                                 | Interceptors for gqlgen                        | NO                         |

## Installation
You can install modules individually, but this can lead to large `go.mod` files.
Instead, we recommend you install a [_recipe_](fx/recipes).

For example:
```shell
go get -v github.com/kevinmichaelchen/fx-libs/fx/recipes/standard
```

### Updating dependencies
If you don't want to use _recipes_, you can update modules together easily:
```shell
go list all | grep github.com/kevinmichaelchen/fx-libs | xargs go get -v
```

## Usage
Check out the [example project](./example), which you can run with
```shell
export NEW_RELIC_KEY=foobar
make example
```

In `main.go`:
```go
package main

import (
	"github.com/kevinmichaelchen/fx-libs/gin"
	"github.com/kevinmichaelchen/fx-libs/recipes/standard"
	"go.uber.org/fx"
)

var Module = fx.Options(
	standard.Module,

	// Gin HTTP handler
	gin.CreateModule(gin.ModuleOptions{
		Invocations: []any{
			// Add a registration function here to register your business logic
			// layer (often called the "Service layer" or "Use case layer") to
			// your GraphQL Handler.
		},
	}),
)

func main() {
	a := fx.New(
		Module,
	)
	a.Run()
}
```

## Contributing

### Releasing new versions
Run the `tag` script by specifying the new version:
```shell
TAG=v0.0.1 ./tag.sh
```

This will generate a Git tag for the repository root, as well as for all descendant submodules.

### `go work sync`
Because this project uses a Go 1.18 Workspace, when working here, you almost
never want to use `go mod`, as it's unable to reason about cross-module
dependencies.

Instead, `go work` is you're friend. When you want to "download or tidy your
dependencies," simply run:
```shell
go work sync
```
