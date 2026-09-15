package handler

import "context"

// NoParent is the typed parent parameter for a root entity's hook state.
// EntityState[T, NoParent].Parent is nil.
type NoParent struct{}

// FieldSet is a read-only view of canonical Go field names. Has returns false
// for unknown fields; it does not resolve JSON names or SQL aliases.
type FieldSet interface {
	Has(field string) bool
}

// OriginalPresence exposes copied pre-initialization marker bits. Available
// distinguishes a missing marker from a known marker whose bits are all false.
// Implementations must not read from the working entity's mutable marker.
type OriginalPresence interface {
	FieldSet
	Available() bool
}

// EntitySnapshot retains an entity's pre-Init/InitMCP processing baseline and
// original presence. It is not the previous database row used for backfill.
// Generated implementations keep captured values, identities and associations
// private and may describe a particular relation path or selected projection.
//
// SyncPresence traverses the captured entity graph, marking changed non-identity
// values while preserving explicit working markers. It reads business values
// without changing them and never mutates the snapshot or reinterprets later IDs
// as original identity. Nil current is a no-op. On error, all working markers
// remain unchanged; implementations validate associations before applying marks.
// Ambiguous entity associations return an error rather than guessing a match.
//
// These are two separate methods: generated owned entities expose
// entity.SyncPresence(snapshot EntitySnapshot[T]) error, which delegates to this
// interface's snapshot.SyncPresence(entity *T) error. Linked entities use the
// latter directly without adding methods to a foreign Go type.
type EntitySnapshot[T any] interface {
	OriginalPresence
	SyncPresence(current *T) error
}

// EntityState contains the typed context of one entity-processing operation.
// Previous must be detached from working values and treated as read-only by
// hooks. PreviousFields identifies loaded fields when previous rows are a
// projection; a missing field must not be interpreted as a known zero value.
// Parent is the declared enclosing relation parent and is nil for a root role.
// SelfParent is the immediate parent through a recursive self edge; it is nil
// at the self-tree root. A self descendant retains the enclosing Parent, so a
// single typed hook contract can observe both contexts without using any.
// Original identity presence is independent
// from any identifier assigned later by initialization or sequencing.
type EntityState[T, P any] struct {
	Previous       *T
	Parent         *P
	Original       OriginalPresence
	PreviousFields FieldSet
	SelfParent     *T
}

// EntityHooks is implemented by a reusable, invocation-scoped hook object.
// Init runs after initial presence synchronization and invariant backfill; it
// uses marker-aware setters when changing business values. Validate then runs
// on the same object and must not mutate values,
// markers or the previous snapshot. Component adapters supply input-specific
// configuration and scoped services without coupling reusable entities to an
// application input type.
type EntityHooks[T, P any] interface {
	Init(context.Context, *T, EntityState[T, P]) error
	Validate(context.Context, *T, EntityState[T, P]) error
}

// AfterSequenceHook is an optional capability of the same invocation-scoped
// EntityHooks object. The generated program calls it after sequencing and before
// diffing. Validated business values and their markers must remain unchanged;
// only identity/link processing allowed by the compiled policy may adjust rows.
type AfterSequenceHook[T, P any] interface {
	AfterSequence(context.Context, *T, EntityState[T, P]) error
}

// AfterQueueHook observes successfully queued mutation work. Queueing is not a
// commit: publishing commit-dependent messages belongs to outcome-aware
// finalization, not this hook. The same object also serves Init and Validate.
type AfterQueueHook[T, P any] interface {
	AfterQueue(context.Context, *T, EntityState[T, P]) error
}

// WriteAction is a planned operation, not evidence that a write was committed.
type WriteAction string

const (
	WriteInsert WriteAction = "insert"
	WriteUpdate WriteAction = "update"
	WriteDelete WriteAction = "delete"
)

// WriteHook optionally customizes business values before validation, using
// marker-aware setters after initial synchronization. Identity reconciliation after diffing is a
// separate phase; this hook is not a commit or message-publication callback.
type WriteHook[T, P any] interface {
	BeforeWrite(context.Context, *T, EntityState[T, P], WriteAction) error
}
