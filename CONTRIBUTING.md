# Contributing

Thanks for helping improve the Datly public SDK. Read [ARCHITECTURE.md](ARCHITECTURE.md)
before changing a contract, and [RELEASING.md](RELEASING.md) for branch selection.
New SDK work targets `v1`; fixes for the existing release line target `main`.

## Develop and verify

Use Go 1.25 or a compatible newer toolchain. Keep the SDK independent of the
Datly runtime and avoid adding dependencies for implementation details.

```sh
go test ./...
go test -race ./...
go vet ./...
```

The dependency-boundary test verifies that SDK packages do not import Datly.
Format edited Go files with `gofmt`. Add tests for observable behavior or contract
invariants when changing executable code. Documentation-only changes should check
links, examples and consistency with the actual exported API.

## Propose a change

Describe the caller's problem, the proposed public behavior and any compatibility
impact. Include relevant test results. Prefer focused interfaces and typed
contracts; database execution, routing, binding engines and transport adapters
belong in runtime or native implementation repositories.

A public interface change can break applications and runtime implementers. Explain
why an additive change is insufficient, update affected examples, and coordinate
matching runtime changes when needed. Do not introduce a dependency from xdatly
back to Datly to make a test or example convenient.

Use a small reproducible example for bug reports. Remove secrets and private
application data. For a potential vulnerability, consult the repository Security
tab for private reporting options; if none is available, request a private
maintainer contact without publishing exploit details or sensitive data.

## Licensing

This repository uses [Apache License 2.0](LICENSE). Preserve existing license and
attribution material when modifying or redistributing source. Identify the origin
and license of any third-party material proposed for inclusion.
