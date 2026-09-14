# Public API by role

Choose the role closest to the work you are doing. The SDK defines public
contracts; the matching Datly runtime provides their implementation.

| Role | What you author | Start with |
| --- | --- | --- |
| Component author | Typed inputs/outputs and component declarations | [Component](component.go), [architecture](ARCHITECTURE.md) |
| Reader author | Row hooks, view selectors, partitioning and reductions | [Reader guide](READERS.md), [lifecycle](handler/lifecycle.go), [selectors](state/selector.go) |
| Cube author | Dimensions, measures, required filters and grouped projections | [Cubes](READERS.md#cube-author) and matching Datly report documentation |
| Cube-composition author | Multiple query frames and their bounded wrapper query | [Composition](READERS.md#cube-composition-author) |
| Handler author | Business operations using invocation capabilities | [Contract](handler/contract.go), [Session](handler/session.go), [Data](handler/data.go) |
| Mutation author | Typed entity hooks, presence, sequencing and completion | [Entity contracts](handler/entity.go), [mutation programs](handler/mutation/program.go) |
| Extension author | Codecs, predicates, docs and optional integrations | [Package responsibilities](ARCHITECTURE.md#package-responsibilities) |
| Runtime implementer | Implementations of contracts, capabilities and lifecycle | [Architecture](ARCHITECTURE.md), [release compatibility](MIGRATION.md) |

A reader does not need a custom write handler. A cube is a reader-derived runtime
feature, not an alternative SDK component or a separate database execution engine.
Roles can overlap within an application; select only the capabilities needed by
the component and keep application authorization explicit.
