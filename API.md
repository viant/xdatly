# Public API Guide

Start with [your role](ROLES.md), or the [reader and cube guide](READERS.md) for
row-hook ordering, view selectors and composition.

This guide covers the exported contract surface in `github.com/viant/xdatly`
for the Datly 1.0 SDK. It is organized by concern rather than package listing:
component authors can start with handlers and sessions, while runtime and
extension authors can jump to providers, codecs, predicates and plugin seams.

The SDK defines interfaces, shared payloads and a few small helpers. Datly
runtime implementations own component discovery, generated handler execution,
binding plans, SQL, transaction orchestration, gateway behavior, authentication
verification, job persistence and provider bootstrapping. Importing a package or
registering a predicate template in this SDK does not automatically expose a
route, install a runtime plugin or register an application component.

For exhaustive field-level documentation, see the
[Go Reference](https://pkg.go.dev/github.com/viant/xdatly). Source links below
point to the files that define each public contract.

## Component And Handler Contracts

The root package exposes [`Component[I, O]`](component.go), a minimal typed
authoring primitive with `Input` and `Output` fields. It describes the shape of a
component without becoming Datly's runtime intermediate representation.

Custom handlers implement [`handler.Contract[I, O]`](handler/contract.go):

```go
package greeting

import (
    "context"

    "github.com/viant/xdatly/handler"
)

type Input struct {
    Name string
}

type Output struct {
    Message string
}

type Handler struct{}

func (*Handler) Exec(ctx context.Context, sess handler.Session, input *Input, output *Output) error {
    output.Message = "Hello, " + input.Name
    return nil
}

var _ handler.Contract[Input, Output] = (*Handler)(nil)
```

[`handler.ContractFunc[I, O]`](handler/contract.go) adapts a function with the
same `Exec` signature. [`handler.Func[I, O]`](handler/func.go) is the pure
function shape for code that only needs `context.Context` and typed input:

```go
package greeting

import (
    "context"

    "github.com/viant/xdatly/handler"
)

type Input struct {
    Name string
}

type Output struct {
    Message string
}

var Pure handler.Func[Input, Output] = func(ctx context.Context, input *Input) (*Output, error) {
    return &Output{Message: "Hello, " + input.Name}, nil
}
```

The SDK does not define route metadata, discovery rules or handler factories.
Those are runtime/build integration concerns.

## Session, Binder And Capabilities

[`handler.Session`](handler/session.go) is intentionally small: it exposes
`Binder() handler.Binder` and `Response() response.Writer`.

[`handler.Binder`](handler/binder.go) has two operations:

| Method | Purpose |
| --- | --- |
| `Bind(ctx, target)` | Populate a destination using the invocation's binding rules. |
| `Lookup(ctx, key)` | Resolve one scoped value. Unknown keys return `(nil, false, nil)`; errors represent failed resolution of a known key. |

Reserved binder keys include `handler.InputKey`, `handler.SelectorsKey`,
`handler.DMLKey`, `handler.SequencerKey`, `handler.FlusherKey`,
`handler.DataKey`, `handler.DifferKey`, `handler.LoggerKey`,
`handler.ValidatorKey`, `handler.MessageBusKey` and `handler.ConnectorKey`.
Additional specialized keys include `handler.InputSnapshotKey`,
`handler.ReadMetadataKey`, `handler.FrameworkValidatorKey`,
`handler.TransactionStarterKey`, `handler.CallerOutputKey` and
`handler.ResultKey`.

Use [`bind.Lookup[T]`](bind/provider.go) and
[`bind.MustLookup[T]`](bind/provider.go) for typed access:

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

The focused data contracts live in [`handler/data.go`](handler/data.go):

| Interface | Methods |
| --- | --- |
| `handler.DML` | `Insert(table, value)`, `Update(table, value)`, `Delete(table, value)`, `Execute(statement, args...)` |
| `handler.MatchedDML` | `UpdateWithOptions(table, value, options...)`, `DeleteWithOptions(table, value, options...)` |
| `handler.Sequencer` | `Allocate(ctx, table, dest, selector)` |
| `handler.Flusher` | `Flush(ctx, table)` |
| `handler.Data` | Embeds `DML`, `Sequencer` and `Flusher` |

`DML` queues write work; sequencing allocates identifiers before queued work is
flushed. `Flush` executes buffered writes. Transaction ownership belongs to the
implementation: SDK contracts distinguish queued work, flushed work and confirmed
commit through completion outcomes.

[`handler.Capabilities`](handler/capabilities.go) groups optional fixed
services: `Differ`, `Logger`, `Validator`, `MessageBus` and `Connector`. `Data`
is kept separate because it is invocation and transaction scoped. The
`Connector` capability is an explicitly granted escape hatch, not managed
transactional data access.

[`bind.Provider`](bind/provider.go) resolves one binder-scoped value by key, and
[`bind.Scope`](bind/provider.go) exposes a binder-backed value scope. These are
contracts for runtime/provider authors; the SDK does not include a reflection
injection engine.

## Input Initialization, Finalization And Outcomes

[`handler/lifecycle.go`](handler/lifecycle.go) defines simple hooks that input,
output or row types may implement:

| Interface | Called for |
| --- | --- |
| `handler.Initializer` | Request-context initialization via `Init(ctx)` |
| `handler.Finalizer` | Success-path follow-up via `Finalize(ctx)` |
| `handler.ErrorFinalizer` | Error-aware completion via `Finalize(ctx, err)` |
| `handler.WriteInitializer` | Mutation entity preparation after keys/sequences and before validation |
| `handler.WriteValidator` | Mutation entity validation before DML is queued |
| `handler.OnFetcher` | Row notification after SQLx population and before relation assembly |
| `handler.OnRelationer` | Row notification after selected child relations are assembled |

[`handler/mcp`](handler/mcp/context.go) adds MCP-specific context and lifecycle
hooks: `mcp.Client`, `mcp.Context`, `mcp.WithContext`,
`mcp.LookupContext`, `mcp.Initializer` and `mcp.Finalizer`. MCP-aware hooks run
only when an MCP context is present; client capabilities are exposed through
`CanElicit`, `CanGenerateContent`, `Elicit` and `GenerateContent`.

[`handler.InputCapturer[I]`](handler/input_capture.go) can capture a detached
input snapshot before ordinary initialization and MCP initialization. The
snapshot is opaque to the SDK and can be exposed through
`handler.InputSnapshotKey`.

[`handler.Outcome`](handler/outcome.go) is the public completion snapshot. It
contains the root operation error plus per-unit `TransactionOutcome` values.
`Outcome.State()` summarizes reported transaction states, and
`Outcome.CommitConfirmed()` returns true only when completion succeeded and every
actual managed transaction committed. Use that method before publishing external
messages that depend on a database commit.

[`handler.WithSession`](handler/session_context.go) and
[`handler.SessionFromContext`](handler/session_context.go) carry a public
session through a context. [`handler.WithHookApplied`](handler/hook_context.go),
[`handler.HookApplied`](handler/hook_context.go),
[`handler.RunPreInvokeHook`](handler/hook_context.go) and
[`handler.InvokeWithSession`](handler/hook_context.go) support pre-invoke hook
coordination for runtime adapters.

[`handler.DataSync`](handler/sync.go) is a small synchronization helper for
cross-view hook ordering. It coordinates readiness by holder key with `Put`,
`Wait`, `Get` and `Delete`; `handler.DataSyncKey` is the context key used by
hook code that shares it.

## Typed Mutation Hook And Policy Seams

Entity-level mutation hooks live in [`handler/entity.go`](handler/entity.go).
They are typed seams for business logic and generated code, not evidence that a
write committed.

| Interface or type | Role |
| --- | --- |
| `handler.FieldSet` | Read-only canonical Go field-name presence checks |
| `handler.OriginalPresence` | Original marker availability and field presence |
| `handler.EntitySnapshot[T]` | Detached pre-initialization snapshot with `SyncPresence(current)` |
| `handler.EntityState[T, P]` | Previous value, parent, self-parent and field evidence for one entity operation |
| `handler.EntityHooks[T, P, O]` | Reusable invocation-scoped `Init` and `Validate` hook object with typed component output access |
| `handler.AfterSequenceHook[T, P, O]` | Optional hook after sequencing and before diffing |
| `handler.AfterQueueHook[T, P, O]` | Optional hook after mutation work is queued |
| `handler.WriteHook[T, P]` | Optional `BeforeWrite` customization with `WriteInsert` or `WriteUpdate` |

`handler.NoParent` is the root parent type. `WriteInsert`, `WriteUpdate` and
`WriteDelete` identify planned actions. Completion is reported separately through
`handler.Outcome`.

Generated writers select `WriteDelete` only from an explicit DQL
`delete_marker(view.column)` and a complete client-supplied identity matched to
an authorized Previous row. Omitting a row or collection does not request deletion.
Unversioned deletion uses `handler.DML.Delete`. A row with a generated
`concurrency_token` requires a client-supplied token and the optional
`handler.MatchedDML` capability. The writer passes `WithIfMatch` so the token
is checked in the SQL `DELETE` itself.

[`handler.Conflict`](handler/conflict.go) carries `Entity`, `Field` and `Reason`
and reports HTTP status 409 through `StatusCode()`, including when wrapped.
Generated `concurrency_token(view.column)` checks compare captured client values
with Previous at the start of validation. The optional `handler.MatchedDML`
capability queues an update or delete with `WithIfMatch(column, previousValue)`.
The SQL statement compares that token in its WHERE clause and requires exactly
one affected row. A later competing write returns `handler.Conflict`; no
database-specific row lock is required. Ordinary `DML` implementations remain
source-compatible; a versioned writer fails closed if `MatchedDML` is absent.
Applications or databases must advance the token on success; the framework
does not invent a next version.

Generated mutation policy contracts live in
[`handler/mutation/program.go`](handler/mutation/program.go):

| Interface | Role |
| --- | --- |
| `mutation.Definition[I, O]` | Immutable generated definition; `Capture` creates one invocation-local program and `FinalizeFailure` handles failures before a program exists |
| `mutation.Program[O]` | Owns invocation-local traversal steps: `Prepare`, `SyncPresence`, `Invariants`, `Init`, `Validate`, `RequiresTransaction`, `Sequence`, `Diff`, `Reconcile`, `Queue`, `Output` and `Finalize` |
| `mutation.Finalizer[I, O]` | Completion callback over input, output and `handler.Outcome` |
| `mutation.AfterSequencer`, `mutation.AfterQueuer` | Optional generated delegation points |

[`handler.Validation`](handler/validation.go) and
[`handler.Violation`](handler/validation.go) define public validation error
payloads. A `Validation` also supplies `StatusCode()`, `ResponseBody()`,
`Err()`, `Error()` and `Messages()`, so response code and body handling can
recognize it through the `response` contracts.

[`handler.ValidationOptions`](handler/validation_options.go) and
`handler.ValidationReference` carry generated mutation policy into the scoped
`handler.Validator` capability. The options describe write action, previous-row
evidence, coverage, deferred fields, satisfied references, connector selection,
location and shallow mode. They carry no database or transaction handles.

[`handler.TransactionStarter`](handler/transaction.go) is an advanced binder
capability keyed by `handler.TransactionStarterKey`. It requests transaction
preparation from the invocation owner; the actual transaction type is opaque.

## Reader Selectors, Read Metadata And Partitions

[`state.Selector`](state/selector.go) represents invocation-scoped query
selection: columns, fields, ordering, offset, limit, page, criteria and
placeholders. `CurrentLimit`, `CurrentOffset`, `CurrentPage`, `SetCriteria` and
`Clone` are helper methods. [`state.NamedSelector`](state/selector.go) binds a
selector to a named view, and [`state.Selectors`](state/selector.go) provides
`Find` and `Clone`.

Selectors are exposed through the binder with `handler.SelectorsKey`. The SDK
does not define a separate reader input channel, and this branch does not export
a type named `reader.QuerySelector`; the public query-selector payload is
`state.Selector`.

Read provenance contracts live in [`handler/read_metadata.go`](handler/read_metadata.go):

| Interface or type | Role |
| --- | --- |
| `handler.ReadMetadata` | Resolves read evidence by canonical input field path |
| `handler.ReadProjection` | Reports field evidence for a completed typed read |
| `handler.ReadOutputProjection` | Resolves projection evidence for output envelope holders |
| `handler.ReadMetadataConsumer` | Opts a handler or definition into read metadata |
| `handler.ReadStep` | Addresses a relation holder and final row ordinal |

Use `handler.WithReadMetadata` and `handler.ReadMetadataFromContext` to carry
read-only evidence through contexts. A nil value intentionally shadows inherited
evidence.

[`reader.Partitioner`](reader/partition.go) and
[`reader.Reducer`](reader/partition.go) are extension contracts for partitioned
typed reads. `PartitionRequest` carries the `*sql.DB`, view name and arguments;
`Partition` supplies a table, expression and args. `Reducer.Reduce` receives the
complete typed result after all batches finish. `ReducerProvider` supplies a
reducer for an invocation.

## Responses, Status And Errors

[`response.Writer`](response/writer.go) is the handler-visible facade for status,
errors and metrics:

```go
package api

import (
    "context"
    "net/http"

    "github.com/viant/xdatly/handler"
)

type Input struct{}
type Output struct {
    Accepted bool
}

func Accept(ctx context.Context, sess handler.Session, input *Input, output *Output) error {
    sess.Response().SetStatusCode(http.StatusAccepted)
    output.Accepted = true
    return nil
}
```

[`response.Response`](response/response.go) is a transport-ready response with
`StatusCode`, `Body`, `Headers`, `Size` and `SetStatusCode`. `response.Compressed`
marks a response whose body is already encoded. [`response.Buffered`](response/buffered.go)
is the SDK's in-memory implementation:

```go
package api

import (
    "net/http"

    "github.com/viant/xdatly/response"
)

func CSVReport() response.Response {
    return response.NewBuffered(
        response.WithStatusCode(http.StatusOK),
        response.WithHeader("Content-Type", "text/csv"),
        response.WithBytes([]byte("id,name\n1,Ada\n")),
    )
}
```

[`response.StatusCoder`](response/status_code.go) lets errors and direct
responses expose status codes. `response.ErrorStatusCode(err, fallback)` follows
ordinary error wrapping and returns the first non-zero status code it finds.

[`response.BodyError`](response/error.go) is an error with an explicit public
body. [`response.Error`](response/error.go) separates public `Code`/`Payload`
from private wrapped `Cause`, and `response.ErrorBody(err)` extracts a public
body through wrapping, including joined errors.

[`response.Status`](response/status.go) is the small public status payload.
[`response.Metric`](response/metric_payload.go), `response.SQLExecution`,
`response.CacheStats`, `response.Metrics`, `response.SQLExecutions` and
`response.ParametrizedSQL` describe execution metrics. Helpers include
`Append`, `Lookup`, `SetError`, `ParametrizedSQL`, `SQL`, `HideSQL`,
`HideMetrics`, `ToSpan`, `ToSpans` and `ExpandSQL`. Those helpers support
logging/debug surfaces; SQL execution itself remains a runtime responsibility.

## Documentation Providers, Codecs And Predicates

[`docs.Source`](docs/source.go) declares an ordered documentation overlay with
global URLs, base URL, one or many document URLs and substitutions. It supports
`Clone`, `Expand`, `IsZero` and `Overlay`.

[`docs.Provider`](docs/service.go) creates a `docs.Service`, and `docs.Service`
looks up payloads by key:

```go
package metadata

import (
    "context"

    "github.com/viant/xdatly/docs"
)

type StaticDocs map[string]string

func (s StaticDocs) Service(ctx context.Context, options ...docs.Option) (docs.Service, error) {
    return s, nil
}

func (s StaticDocs) Lookup(ctx context.Context, key string) (string, bool, error) {
    value, ok := s[key]
    return value, ok, nil
}

var _ docs.Provider = StaticDocs{}
var _ docs.Service = StaticDocs{}
```

[`docs.Options`](docs/options.go), `docs.WithURL` and `docs.WithConnector`
configure providers. `docs.Connector` is a bounded compatibility seam exposing
`DB() (*sql.DB, error)`; it does not make documentation ownership a database
runtime concern.

[`codec.Config`](codec/codec.go) describes a codec's source and destination
types, body, args and output expression. `codec.Factory` creates instances, and
`codec.Instance.Value(ctx, raw, options...)` performs conversion. Public options
in [`codec/options.go`](codec/options.go) include type lookup, columns source,
selector, record, value lookup, value getter, resource filesystem and extra
options. Registered codec instances may be shared by concurrent invocations, so
implementations must not retain invocation-local pointers or context.

[`predicate.Template`](predicate/template.go) and `predicate.NamedArgument`
describe predicate SQL templates. `predicate.Lookup` is the read-only registry
shape. [`predicate.Registry`](predicate/registry.go) supports `Lookup`, while
package helpers `predicate.New`, `predicate.RegisterTemplate` and
`predicate.Templates` manage the package-level template registry and optional
notifications. This is a template registry only; it does not expose components.

[`predicate.Handler`](predicate/handler.go) computes one SQL criteria fragment
from a bound value, and `predicate.HandlerFunc` adapts a function. Filter
payloads in [`predicate/filter.go`](predicate/filter.go) include
`NamedFilters`, `NamedFilter`, `StringsFilter`, `IntFilter` and `BoolFilter`.

## Async Job Contracts

[`async.Job`](async/job.go) is the public async job metadata payload. It embeds
destination table settings from [`async/destination.Table`](async/destination/table.go),
cache settings from [`async/destination.Cache`](async/destination/cache.go),
request metadata, principal metadata, timings, status and error text. It is a
payload contract, not the job persistence engine.

Public constants include [`async.StatusPending`, `StatusRunning`, `StatusDone`
and `StatusError`](async/status.go), [`async.InvocationTypeEvent` and
`InvocationTypeUndefined`](async/invocation.go), and notification methods in
[`async/notification.go`](async/notification.go).

Destination dispositions live in
[`async/destination/disposition.go`](async/destination/disposition.go):
`CreateDispositionIfNeeded`, `CreateDispositionNever`,
`WriteDispositionTruncate`, `WriteDispositionEmpty` and
`WriteDispositionAppend`.

## Execution, Tracing, Logging And Message Bus

[`exec.Context`](exec/context.go) is the public execution-context payload for
method, URI, status, errors, elapsed time, headers, metrics and trace data. Use
`exec.New(...)` with options such as `WithMethod`, `WithURI`, `WithStatusCode`,
`WithStatus`, `WithStartTime`, `WithTrace`, `WithTraceResource` and
`WithHeaders`. `exec.WithContext` and `exec.GetContext` attach and retrieve it
from `context.Context`.

The context also exposes `AppendMetrics`, `SetError`, `SetValue`, `Value`,
`SnapshotForLogging` and `Complete`. [`exec.NewHTTPContext`](exec/context_http.go)
builds an HTTP-colored execution context, and `exec.NewContext` is the
compatibility wrapper for the older HTTP constructor shape. [`exec.PrepareLogging`](exec/logging.go)
returns an audit snapshot and optional trace, hiding SQL when requested.

[`tracing.Trace`](tracing/trace.go), `tracing.Span`, `tracing.SpanStatus` and
`tracing.ResourceInfo` are reduced trace payloads. Helpers include
`tracing.NewTrace`, `tracing.NewSpan`, `Trace.Append`, `Span.OnDone`,
`Span.SetStatus`, `Span.SetStatusFromHTTPCode` and `Span.WithAttributes`.
The constants are `tracing.StatusOK`, `tracing.StatusError` and
`tracing.StatusUnset`.

[`logger.Logger`](logger/logger.go) is the minimal per-invocation logging
interface: `Debug`, `Info`, `Warn` and `Error`.

[`mbus.Service`](mbus/service.go) creates and pushes messages. `mbus.Message`
contains ID, resource, trace ID, attributes, subject and data; `mbus.Confirmation`
contains the delivered message ID. Options include `WithAttribute`,
`WithTraceID`, `WithID` and `WithSubject`.

```go
package events

import (
    "context"

    "github.com/viant/xdatly/handler"
)

func Publish(ctx context.Context, bus handler.MessageBus, outcome handler.Outcome, orderID int) error {
    if !outcome.CommitConfirmed() {
        return nil
    }
    message := bus.Message("orders.created", map[string]int{"id": orderID})
    _, err := bus.Push(ctx, message)
    return err
}
```

## Plugin, Connector, Differ And Auth

[`plugin.Registry`](plugin/registry.go) is a declarative registration seam for
the contract families in this SDK: codecs, codec factories, predicate templates
and docs providers. It is not a hidden global plugin bootstrapper and does not
replace runtime component discovery.

[`connector.Provider`](connector/provider.go) resolves an exact registered
connector name to a borrowed `*sql.DB`. Empty or unknown names and canceled
contexts should fail. Direct use of the returned DB is outside Datly's managed
`handler.Data` transaction and flush lifecycle, and handlers must not close the
borrowed DB.

[`differ.Differ`](differ/differ.go) compares two values without mutating them
and returns a [`differ.ChangeLog`](differ/log.go). Public change types are
`ChangeTypeCreate`, `ChangeTypeUpdate` and `ChangeTypeDelete`.
[`differ.Path`](differ/path.go) builds field, map-entry and slice-index paths;
`ChangeLog.ToChangeRecords` adds optional source/source ID/user ID metadata
through `WithSource`, `WithSourceID` and `WithUserID`. `ChangeLog.Err()` joins
native comparison failures retained on change records.

[`auth.Authenticator`](auth/auth.go) authenticates an extracted `auth.Token` and
returns an `auth.Principal` with normalized `auth.Claims`. Vendors and common
claim names are constants: `VendorJWT`, `VendorCognito`, `VendorFirebase`,
`ClaimSubject`, `ClaimEmail` and `ClaimScope`. Authentication errors include
`ErrInvalidToken` and `ErrExpiredToken`. Transport extraction and vendor-specific
verification belong to implementations.

## Low-Level Types And Reference

This guide expands the public interfaces and the payload types most applications
touch directly. Several exported structs are intentionally simple data carriers:
metrics, async destination fields, validation references, trace resource fields
and filter payloads. Their exact fields and JSON tags are best read in the
[Go Reference](https://pkg.go.dev/github.com/viant/xdatly) or the linked source
files when implementing serialization, compatibility checks or provider tests.
