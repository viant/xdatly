package handler

// FrameworkValidatorKey is reserved for the mandatory generated validation
// phase. It uses the same Validator interface, but cannot be replaced by request
// values or the optional custom-handler Validator capability.
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
	Action WriteAction
	// Previous is the detached typed database row matched using original
	// identity. Inserts require nil; updates require a row with its complete
	// mapped identity. Unique checks exclude that tuple, never working IDs.
	Previous any
	// PreviousFields is actual read provenance for Previous. Updates require
	// it; unknown identity or constraint dependencies cannot be replaced by
	// zero values from a partial projection. Inserts require nil.
	PreviousFields FieldSet
	// Fields is effective validation coverage, not the persistence Has mask.
	// Inserts require nil (complete checks). Updates require explicit coverage;
	// invariant-backfilled fields may be covered without becoming client writes.
	// Names are canonical Go field names, as with EntityState.PreviousFields.
	Fields FieldSet
	// DeferredFields identifies unavailable INSERT inputs for a generated
	// business-validation pass. It is not a sparse write mask. Only compiled
	// producer authority and original absence may supply this set; a mandatory
	// final pass must validate the completed payload with no fields deferred.
	// UPDATE and calls with SatisfiedReferences cannot defer fields.
	DeferredFields FieldSet
	// SatisfiedReferences identifies exact reference-existence checks already
	// proven by the generated pending-write graph for an INSERT payload. Other
	// constraints remain active. The generated owner must prove matching typed
	// parent keys, INSERT action and write order in the same transaction scope.
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
