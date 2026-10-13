package handler

import "context"

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
