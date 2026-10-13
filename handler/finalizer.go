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
