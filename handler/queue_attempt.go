package handler

import "context"

// QueueDisposition distinguishes physical work from a matched identity-only
// update. Skipped deletes and unmatched identities are not queue attempts.
type QueueDisposition string

const (
	QueuePhysical QueueDisposition = "physical"
	QueueNoop     QueueDisposition = "noop"
)

type QueueResult string

const (
	QueueQueued        QueueResult = "queued"
	QueueNoopCompleted QueueResult = "noop"
	QueueFailed        QueueResult = "failed"
	QueueCanceled      QueueResult = "canceled"
	QueuePanicked      QueueResult = "panicked"
)

// QueueAttemptEvent describes actual queue traversal, never transaction
// completion. Row and Previous are detached pointers to the canonical role's
// entity type. Presence uses canonical Go field names. Each boundary owns fresh
// evidence; changing it cannot change persistence or another callback's evidence.
// EvidenceError reports unavailable evidence without changing execution. Cause
// and EvidenceError are private diagnostics, not public response messages.
// Cause supports errors.Is against the retained operational cause. It does not
// expose mutable typed causes through errors.As/Unwrap; execution owns them.
type QueueAttemptEvent struct {
	InvocationID uint64
	Attempt      int
	Position     int
	Location     string
	Role         string
	Operation    WriteAction
	Disposition  QueueDisposition
	Boundary     PhaseBoundary
	Result       QueueResult
	// Queued is true once DML returned success, even if AfterQueue failed.
	// False means queue acceptance was not confirmed.
	Queued         bool
	Row            any
	Previous       any
	Presence       FieldSet
	Original       OriginalPresence
	PreviousFields FieldSet
	Cause          error
	EvidenceError  error
}

// QueueAttemptObserver is opt-in on a physical writer root's lifecycle object.
// It receives all graph roles synchronously. It has no output or persistence
// capability and cannot veto work. Implementations must not retain invocation
// services or use them to alter state. Panics are contained by the runtime.
type QueueAttemptObserver interface {
	ObserveQueueAttempt(context.Context, QueueAttemptEvent)
}
