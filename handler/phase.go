package handler

import "context"

// InvocationPhase identifies a real execution boundary, not a transaction result.
type InvocationPhase string

const (
	PhaseInvocation          InvocationPhase = "invocation"
	PhaseExecution           InvocationPhase = "execution"
	PhaseInputInitialization InvocationPhase = "input_initialization"
	PhaseBinding             InvocationPhase = "binding"
	PhaseInitialization      InvocationPhase = "initialization"
	PhaseValidation          InvocationPhase = "validation"
	PhaseAllocation          InvocationPhase = "allocation"
	PhaseQueue               InvocationPhase = "queue"
)

type PhaseBoundary string

const (
	PhaseBegin PhaseBoundary = "begin"
	PhaseEnd   PhaseBoundary = "end"
)

type PhaseResult string

const (
	PhaseSucceeded  PhaseResult = "succeeded"
	PhaseViolations PhaseResult = "violations"
	PhaseFailed     PhaseResult = "failed"
	PhaseCanceled   PhaseResult = "canceled"
	PhasePanicked   PhaseResult = "panicked"
)

// PhaseEvent carries private diagnostics. Cause must not be exposed to clients
// or logged without the application's existing diagnostic policy. A begin event
// has no result. Outcome remains the sole authority for transaction completion.
type PhaseEvent struct {
	Phase        InvocationPhase
	Boundary     PhaseBoundary
	Result       PhaseResult
	InvocationID uint64
	Attempt      int
	Cause        error
}

// PhaseObserver observes execution without changing its result. Before binding,
// it must use only explicitly available early-safe capabilities (the context
// logger), never assume a bound input, claims, or Current data. Implementations
// must not mutate business state or retain invocation objects across requests.
type PhaseObserver interface {
	ObservePhase(context.Context, PhaseEvent)
}
