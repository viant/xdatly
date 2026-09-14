# Migrating to the v1 SDK

The `v1` branch is a new public contract surface. It keeps the canonical module
path `github.com/viant/xdatly`, but replaces the previous nested-module layout
with one module. It is not a drop-in update for every previous handler or session.
Keep existing applications on their current pinned versions until migrated.

## Coordinate module versions

Use matching Datly and xdatly v1 revisions. After the SDK branch is published,
`go get github.com/viant/xdatly@v1` resolves the branch to a concrete module version;
commit that resolved version for reproducibility. A branch lookup is not a
`v1.0.0` release tag. See [RELEASING.md](RELEASING.md).

Remove obsolete requirements on separately versioned SDK submodules after updating
imports, then run `go mod tidy` and the application's tests. Do not keep old and
new module providers for the same import path: they can cause ambiguous imports.

## Update application contracts

| Previous package area | v1 package area |
| --- | --- |
| `handler/async` | `async` |
| `handler/auth` | `auth` |
| `handler/exec` | `exec` |
| `handler/response` | `response` |
| `handler/state` | `state` |
| `handler/tracing` | `tracing` |
| `handler/logger` | `logger` |
| `handler/mbus` | `mbus` |
| `handler/differ` | `differ` |

These are responsibility mappings, not promises of identical types or signatures.
Review each contract at the new package. The old `types/core`, `types/custom` and
`extension` nested modules are not part of this SDK.

Implement `handler.Contract[I, O]` for typed handlers. The reduced session exposes
`Binder()` and `Response()`. Request scoped `DML`, `Sequencer`, `Flusher` or `Data`
capabilities instead of depending on the previous broad session surface. Runtime
implementations and component discovery belong to Datly.

Preserve transaction semantics when migrating: queued DML is not a confirmed
commit, and flushing a caller-owned transaction does not transfer commit ownership.
Use completion outcomes for messages that depend on a successful commit.

## Verify the migrated application

Compile the application against the matched SDK/runtime revisions. Exercise real
input binding, handler execution, output/error handling and transaction completion.
For exposed components, verify HTTP/MCP behavior as well as compilation. Check
initialization/finalization hooks, sparse inputs, null values and any async replay
or message publication the application uses.
