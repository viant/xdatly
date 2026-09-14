package handler

import "context"

// Initializer is the reduced public pre-execution lifecycle hook for inputs or
// other handler-bound values that need request-context initialization without
// binding themselves to Datly runtime internals.
type Initializer interface {
	Init(ctx context.Context) error
}

// Finalizer is the reduced public success-path lifecycle hook for outputs or
// other handler-bound values that need post-execution cleanup/follow-up without
// exposing transport- or runtime-specific finalization concerns.
type Finalizer interface {
	Finalize(ctx context.Context) error
}

// ErrorFinalizer is the reduced public error-aware lifecycle hook for outputs
// that need the terminal handler error when finalizing.
type ErrorFinalizer interface {
	Finalize(ctx context.Context, err error) error
}

// WriteInitializer is implemented by typed entities that need preparation
// after generated relation keys and sequences are available but before a
// generated handler validates or queues the entity for persistence.
type WriteInitializer interface {
	InitWrite(ctx context.Context) error
}

// WriteValidator is implemented by typed entities that need validation after
// write initialization and before a generated handler queues DML.
type WriteValidator interface {
	ValidateWrite(ctx context.Context) error
}

// OnFetcher is implemented by typed rows that need initialization immediately
// after SQLx has populated the row and before Datly indexes or relates it.
type OnFetcher interface {
	OnFetch(ctx context.Context) error
}

// OnRelationer is implemented by typed rows that need notification after all
// selected child relations have been assembled.
type OnRelationer interface {
	OnRelation(ctx context.Context)
}
