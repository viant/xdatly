package handler

// FrameworkValidatorKey is the protected native validation capability used by
// generated writers and trusted application-owned buffered writers. It cannot
// be replaced by request values or the optional custom-handler Validator.
const FrameworkValidatorKey ValueKey = "frameworkValidator"

// ValidationOptions supplies mutation policy to the scoped Validator capability.
// It carries no database or transaction handles. The capability rejects unknown
// connectors, unknown actions and inconsistent options before performing checks.
// Generated entity validation is shallow: all frames are validated by the
// program before any custom EntityHooks.Validate method runs.
// The framework capability accepts either one entity and one ValidationOptions,
// or a typed entity slice and one []ValidationOptions argument of equal length.
// Batch policies retain each candidate's matched prior row and diagnostic
// location; candidates in a batch must share the entity type and connector.
type ValidationOptions struct {
	// CheckOnly runs selected Go, not-null and reference constraints without
	// mutation provenance. It requires empty Action,
	// explicit Fields, Shallow and CheckUnique=false. Mutation provenance,
	// HonorPresence, deferrals and satisfied references are forbidden. It uses
	// the scoped connection/transaction without flushing queued DML. A typed
	// batch must use CheckOnly for every candidate; Go checks precede SQL checks.
	CheckOnly bool
	Action    WriteAction
	// CheckUnique and CheckRef select native database checks for this pass.
	// Nil preserves the default true. Required and Go checks remain active
	// under Fields/DeferredFields. Generated final passes always enable both.
	CheckUnique *bool
	CheckRef    *bool
	// Previous is the detached typed database row matched using original
	// identity. Inserts require nil; updates require a row with its complete
	// mapped identity. Unique checks exclude that tuple, never working IDs.
	Previous any
	// PreviousFields is actual read provenance for Previous. Updates require
	// it; unknown identity or constraint dependencies cannot be replaced by
	// zero values from a partial projection. Inserts require nil.
	PreviousFields FieldSet
	// Fields is effective validation coverage, not the persistence Has mask.
	// Inserts require nil unless HonorPresence is enabled. Updates require explicit coverage;
	// invariant-backfilled fields may be covered without becoming client writes.
	// Names are canonical Go field names, as with EntityState.PreviousFields.
	Fields FieldSet
	// GoFields optionally selects independent Go-check coverage for CheckOnly.
	// Nil (including typed nil) defaults to Fields, never to all fields. An
	// explicit all-fields set requests complete Go checks. Both sets are
	// snapshotted before connection access or predicate execution. Other modes
	// reject effective nonnil GoFields.
	GoFields FieldSet
	// HonorPresence enables supplied-field coverage for an INSERT business pass.
	// It requires Fields and forbids Previous evidence. Generated writers always
	// run a complete final validation pass before queuing the write.
	HonorPresence bool
	// DeferredFields identifies unavailable INSERT inputs for a generated
	// business-validation pass. It is not a sparse write mask. Only compiled
	// producer authority and original absence may supply this set; a mandatory
	// final pass must validate the completed payload with no fields deferred.
	// UPDATE and calls with SatisfiedReferences cannot defer fields.
	DeferredFields FieldSet
	// SatisfiedReferences identifies exact reference-existence checks already
	// proven for an INSERT payload. Other constraints remain active. A generated
	// owner or trusted application-owned buffered writer must successfully queue
	// the exact typed parent INSERT through the same invocation DML owner, prove
	// its earlier order and table, match captured parent/child keys through native
	// reference matching, and preserve those bindings until completion.
	// These are trusted framework facts, never request-bindable assertions.
	SatisfiedReferences []ValidationReference
	// Connector selects a configured connector within the invocation scope.
	// Empty selects that scope's configured default, not an arbitrary database.
	Connector string
	Location  string
	Shallow   bool
}

// ValidationReference identifies one native reference constraint exactly.
// Empty schema remains empty; field names are canonical candidate Go fields.
type ValidationReference struct {
	Field  string
	Schema string
	Table  string
	Column string
}
