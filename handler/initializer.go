package handler

import "context"

// Initializer is the reduced public pre-execution lifecycle hook for inputs or
// other handler-bound values that need request-context initialization without
// binding themselves to Datly runtime internals.
type Initializer interface {
	Init(ctx context.Context) error
}

// InvokerInitializer prepares a bound input before handler execution. Prefer it
// when initialization composes other components instead of storing an invoker
// on the input. The invoker is valid only during Init; child calls must be
// awaited and retain the caller's providers and transaction owner.
type InvokerInitializer interface {
	Init(ctx context.Context, invoker ComponentInvoker) error
}
