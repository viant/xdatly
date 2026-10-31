package handler

import "context"

// InvokerFinalizer composes dependencies after a handler returns an output and
// before transaction completion. Prefer it over an injected output/input
// invoker. The invoker is valid only during Finalize and calls must be awaited.
// Hook and child errors prevent owned commit; an existing handler error is
// preserved. It does not observe early binding/initialization failures or a
// committed outcome. Outcome finalization remains a separate phase after
// transaction completion.
type InvokerFinalizer interface {
	Finalize(ctx context.Context, invoker ComponentInvoker) error
}

// InvokerErrorFinalizer composes dependencies before transaction completion
// with the accumulated operation error at the callback boundary: handler,
// output-binding, or preparation failures encountered before this call. The
// error is nil on success and is retained even when Finalize returns nil.
// The invoker is valid only during Finalize; child calls must be awaited.
// Hook and child errors prevent owned commit. This hook does not observe early
// failures that produced no output, or future drain/commit failures. Outcome
// finalization remains a separate phase after transaction completion.
type InvokerErrorFinalizer interface {
	Finalize(ctx context.Context, invoker ComponentInvoker, operationErr error) error
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

// BeforeCommitValidator is an optional output validation hook. The native owner
// calls it once after every buffered write has physically drained and before any
// locally owned transaction commits. A returned error prevents owned commit.
// Earlier operation/drain errors skip this hook; Finalize and outcome-aware
// finalization retain their existing phases.
//
// The supplied invoker is valid only during this callback and admits canonical
// native readers over the owner's already-open transactions. It cannot enqueue
// writes, open independent units, publish messages, stream rows, or complete the
// owner. Calls must be awaited. Managed violations remain failures even if the
// hook catches their immediate error. Do not retain the invoker or launch work
// beyond the callback. Captured writer values and presence remain read-only.
//
// Callback context retains caller identity, values, deadline and cancellation,
// but not the owner's database binding. A separately owned external authority
// reader therefore retains its independent committed view; only this invoker
// receives the pending owner transaction. Caller-supplied transactions retain
// caller completion ownership: validation precedes handoff, not a later commit
// performed outside the native invocation. This contract does not sandbox
// arbitrary Go handles or external side effects.
type BeforeCommitValidator interface {
	BeforeCommit(context.Context, ComponentInvoker) error
}
