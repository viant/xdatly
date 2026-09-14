# xdatly

[![Go Reference](https://pkg.go.dev/badge/github.com/viant/xdatly.svg)](https://pkg.go.dev/github.com/viant/xdatly)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)](go.mod)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[Package documentation](https://pkg.go.dev/github.com/viant/xdatly) ·
[Public API guide](API.md) ·
[Architecture](ARCHITECTURE.md) · [Migration](MIGRATION.md) ·
[Contributing](CONTRIBUTING.md) · [Releases](RELEASING.md)

`xdatly` is the public Go SDK for Datly 1.0 contracts. It gives application
authors, runtime implementers and extension providers one small module for typed
components, handlers, invocation capabilities, response payloads, documentation
providers, codecs, predicates, async job metadata and optional integration seams.

This branch uses one Go module, `github.com/viant/xdatly`, with Go 1.25. The
SDK contains contracts and small shared primitives; `github.com/viant/datly`
supplies the runtime implementation for routing, binding, SQL execution,
transactions, gateways and generated handlers. The SDK does not import Datly and
does not automatically register components or providers.

## Install

Use the canonical module path:

```sh
go get github.com/viant/xdatly@v1
```

Pin the resolved version with the matching Datly runtime revision. A branch named
`v1` does not add `/v1` to Go imports and does not imply a `v1.0.0` release tag.
See [RELEASING.md](RELEASING.md).

## Typed handler example

```go
package greeting

import (
    "context"
    "github.com/viant/xdatly/handler"
)

type Input struct { Name string }
type Output struct { Message string }
type Handler struct{}

func (*Handler) Exec(ctx context.Context, session handler.Session, input *Input, output *Output) error {
    output.Message = "Hello, " + input.Name
    return nil
}

var _ handler.Contract[Input, Output] = (*Handler)(nil)
```

The runtime supplies the bound input, invocation session and output. Route
metadata, component discovery and linking are configured in Datly; this example
defines the SDK contract only.

## Capability example

Use `handler.Session.Binder()` to resolve invocation-scoped capabilities. The
same key names are public contracts; the concrete services are supplied by the
runtime or by explicitly registered providers.

```go
package orders

import (
    "context"
    "fmt"

    "github.com/viant/xdatly/bind"
    "github.com/viant/xdatly/handler"
)

type CreateOrderInput struct {
    CustomerID int
    Amount     float64
}

type CreateOrderOutput struct {
    ID int
}

func Create(ctx context.Context, sess handler.Session, input *CreateOrderInput, output *CreateOrderOutput) error {
    data, ok, err := bind.Lookup[handler.Data](ctx, sess.Binder(), handler.DataKey)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("data capability is not available")
    }
    if err := data.Allocate(ctx, "orders", &output.ID, "id"); err != nil {
        return err
    }
    if err := data.Insert("orders", input); err != nil {
        return err
    }
    return data.Flush(ctx, "orders")
}
```

## Public packages

| Package | Public surface |
| --- | --- |
| `github.com/viant/xdatly` | `Component[I, O]` typed authoring primitive |
| `handler` | `Contract`, `ContractFunc`, `Session`, `Binder`, data capabilities, lifecycle hooks, validation, outcomes and read metadata |
| `handler/mutation` | Generated mutation definition, program and completion contracts |
| `handler/mcp` | MCP client context plus MCP-aware input/output lifecycle hooks |
| `bind`, `state`, `reader` | Binder providers and typed lookup, selectors and read partition/reducer contracts |
| `response` | Handler-visible response writer, direct transport responses, status/errors and metrics |
| `docs`, `codec`, `predicate` | Documentation lookup, value conversion and predicate extension contracts |
| `async`, `exec`, `tracing`, `logger`, `mbus` | Async job payloads, execution context, trace payloads, logging and messaging contracts |
| `plugin`, `connector`, `differ`, `auth` | Registration, direct connector access, comparison/change logs and authentication contracts |

Read [API.md](API.md) for a guided pass through the public interfaces, including
source links and practical examples. Low-level payload fields and helper methods
that are not expanded there are available in the
[Go Reference](https://pkg.go.dev/github.com/viant/xdatly).

Use `github.com/viant/datly` for runtime implementation and authoring guides.

See [ARCHITECTURE.md](ARCHITECTURE.md) for dependency direction, invocation,
capability ownership and package boundaries.

## Development and releases

The `main` branch retains the existing SDK release line. This `v1` branch contains
the new single-module SDK. It is under release validation; creating or pushing
the branch does not publish a `v1.0.0` release. See [RELEASING.md](RELEASING.md).

Run the package and dependency-boundary checks with:

```sh
go test ./...
```

## Documentation and contributing

- [Architecture](ARCHITECTURE.md): dependency direction, invocation and ownership.
- [Public API guide](API.md): public contracts grouped by component concern.
- [Migration](MIGRATION.md): moving from the previous SDK layout to v1.
- [Contributing](CONTRIBUTING.md): development checks and change guidance.
- [Releases](RELEASING.md): branches, module identity and publication order.

Report reproducible bugs and feature requests in the repository issue tracker.
Include the SDK revision, Go version and a minimal example; remove credentials
and private application data before sharing diagnostics.

## License

Licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE) for the
full terms. The original repository license is retained unchanged. Dependencies
retain their own licenses and notices.
