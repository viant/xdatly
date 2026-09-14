# Reader, cube and composition authors

[Roles](ROLES.md) · [Architecture](ARCHITECTURE.md)

## Reader author

Declare typed input and output contracts using `Component[I, O]`. Datly's reader
runtime interprets the authored view/SQL metadata, binds rows and assembles
relations. The SDK supplies the query-selection, hook and extension contracts.

### Row hooks and ordering

[handler/lifecycle.go](handler/lifecycle.go) defines:

```go
type OnFetcher interface {
    OnFetch(ctx context.Context) error
}
type OnRelationer interface {
    OnRelation(ctx context.Context)
}
```

`OnFetch` runs after SQLX populates a row and before Datly indexes or relates it.
Return an error to fail that fetch. `OnRelation` runs after the selected child
relations have been assembled; it does not return an error. Do not assume omitted
relations were loaded, or use this notification to report a validation failure.

```go
package rows

import (
    "context"
    "strings"
    "github.com/viant/xdatly/handler"
)

type Product struct { ID int }
type Inventory struct {
    ID int
    Name string
    Products []*Product
    Label string `sqlx:"-"`
    SelectedProductCount int `sqlx:"-"`
}

func (r *Inventory) OnFetch(ctx context.Context) error {
    r.Label = strings.TrimSpace(r.Name)
    return nil
}
func (r *Inventory) OnRelation(ctx context.Context) {
    r.SelectedProductCount = len(r.Products)
}

var _ handler.OnFetcher = (*Inventory)(nil)
var _ handler.OnRelationer = (*Inventory)(nil)
```

The count describes the selected relation result, not the total number of related
rows in the database. Keep hook work appropriate to a read path; a hook may also
be exercised by nested reads, warmup or retries according to runtime policy.
Do not make externally visible publication depend on counting hook invocations.

Input `Initializer.Init(context.Context) error` and output finalization are
invocation lifecycle hooks, separate from per-row hooks. Error-aware and
injector-aware finalization contracts are defined in the handler package. Choose
the contract supported by the runtime invocation instead of manually calling
finalizers from a row hook.

### Selectors and read evidence

[state.Selector](state/selector.go) carries fields/columns, ordering, offset,
limit, page, criteria and placeholders. `NamedSelector` addresses a named
view/query scope. Datly validates requested names against the authored source
projection and permitted selector metadata. Main-view defaults do not implicitly
apply to every subview: bind selectors to their corresponding view.

Requesting only root inventory fields can omit its product relation. A separately
bound product selector controls the selected product fields, filters and paging.
Author the allowed output names explicitly; the runtime must not invent join
aliases or silently accept unrelated columns.

[ReadMetadata](handler/read_metadata.go) describes which fields were actually
loaded for bound inputs. `ReadProjection.Fields` addresses root ordinals and
optional `ReadStep` relation paths. Unknown evidence is an error; it must not be
treated as proof that all fields were loaded. `ReadOutputProjection` addresses
additional output-envelope slots.

### Partitions and reducers

[reader.Partitioner](reader/partition.go) returns the authored partitions for a
read. `ReducerProvider` supplies a reducer, whose `Reduce` receives the complete
typed result for the view. Results across partitions/batches have no guaranteed
input order; sort explicitly for an order-sensitive reduction. Relation readers
invoke reduction after all batches finish. Preserve the expected result shape.

## Cube author

Cubes add dimensions, measures and aggregation semantics to a declared reader.
Their metadata, compilation and registration APIs belong to Datly's `report`
packages, not to an SDK `Cube` interface. The same input authorization and
required filters must remain effective in derived reads.

Document which dimensions determine grouping and which measures are aggregates.
Selecting fewer measures can reuse a warmed full projection while retaining all
warmed grouping dimensions. Dropping a dimension changes grouping and can require
reaggregation; simply hiding that column is not a valid grouped-cache answer.
Ordinary readers can narrow columns without this cube grouping requirement.

## Cube-composition author

Composition prepares multiple query frames and combines them with a bounded
wrapper query. Datly's report composition supports `$CubeSQL1` through
`$CubeSQLN`; `inheritFrom` refers to a prior frame using one-based numbering.
Use declared contract columns and explicit output aliases, preserving each
frame's required filters, authorization and ordered SQL arguments.

For example, two frames can compare web and store spend for the same tenant and
account. The second frame inherits common filters and overrides the channel.
The wrapper joins the two declared frame results and explicitly names the web
and store measure outputs. It does not accept arbitrary database tables as a
substitute for the registered cube sources.

Configure cube-count, result and timeout budgets in the matching Datly report
API. Composition must also respect the database's combined placeholder limit.
The current standalone source snapshot requires report registration through an
embedding host; report metadata alone does not activate derived endpoints in
the standalone loader. Check the matching Datly release's report guide before
using this feature in a deployment.

## Verify behavior

Exercise row hooks before/after relation assembly, omitted child relations,
multiple subview selectors, pagination and filters. Test cache misses and hits,
regular projection narrowing and cube measure narrowing separately. For
composition, cover inherited filters, nullable joins, explicit aliases, denied
rows and exceeded budgets through the transports the application exposes.
