package handler

import "context"

const (
	// DMLKey identifies the invocation-scoped buffered write capability.
	DMLKey ValueKey = "dml"
	// SequencerKey identifies the invocation-scoped identifier allocator.
	SequencerKey ValueKey = "sequencer"
	// FlusherKey identifies explicit execution of buffered writes.
	FlusherKey ValueKey = "flusher"
	// DataKey identifies the complete invocation data capability.
	DataKey ValueKey = "data"
)

// DML is the narrow buffered write contract available to handlers.
// Transaction completion is available separately through Flusher and Data.
type DML interface {
	Insert(table string, value any) error
	Update(table string, value any) error
	Delete(table string, value any) error
	Execute(statement string, args ...any) error
}

// Sequencer allocates identifiers before buffered writes are flushed.
type Sequencer interface {
	Allocate(ctx context.Context, table string, dest any, selector string) error
}

// Flusher executes buffered writes. Implementations complete transactions they
// open themselves and leave supplied transactions to their owner.
type Flusher interface {
	Flush(ctx context.Context, table string) error
}

// Data is the complete handler data capability. Focused consumers may request
// DML, Sequencer, or Flusher independently.
type Data interface {
	DML
	Sequencer
	Flusher
}
