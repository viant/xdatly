# xdatly

Public Go contracts for Datly 1.0: build typed components, implement handlers and
share the same application contracts across HTTP and MCP.

This branch uses one Go module, `github.com/viant/xdatly`, with Go 1.25.
The SDK contains contracts and small shared primitives; the Datly runtime
implements execution, binding, SQL, transactions and gateway behavior. The SDK
does not depend on the Datly implementation.

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

The runtime supplies the bound input, invocation session and output. Route metadata,
component discovery and linking are configured in Datly; this example defines the
SDK contract only.

## Public packages

- Root `Component[I, O]` declares typed components.
- `handler` and `handler/mutation` define handlers, scoped capabilities, hooks and completion.
- `bind`, `state`, `codec` and `predicate` describe application inputs and extensions.
- `reader`, `response` and `docs` define query selection, output and documentation contracts.
- `async`, `exec`, `tracing`, `logger` and `mbus` expose execution capabilities.
- `plugin`, `connector` and `differ` define optional integrations.

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
