package handler

import "context"

// Route selects a component by its method and concrete URL (including query
// parameters). Scope and Name optionally require an exact registered component
// identity; route lookup still enforces the runtime's exposed route catalog.
type Route struct {
	Method string
	URL    string
	Scope  string
	Name   string
}

// InjectorLookup resolves a route without executing it. The returned Binder
// executes that component once on its first Bind or ResultKey lookup. Bind uses
// native binding tags to match the component result to the destination.
// Both lookup and binders are valid only during the Finalize call; calls must be
// awaited. Call contexts supply cancellation while retaining invocation values.
type InjectorLookup = func(context.Context, Route) (Binder, error)

const (
	// CallerOutputKey is the declared source for a conditional child's input fields.
	// Only fields authored with this source receive the calling output's values.
	CallerOutputKey ValueKey = "caller_output"
	// ResultKey resolves the complete typed result of a route-specific binder.
	ResultKey ValueKey = "result"
)

// InjectorFinalizer composes conditional dependencies before transaction
// completion. Its child work participates in the root unit of work. Errors
// prevent owned commit; an existing operation error is never discarded.
// It takes precedence over ordinary output finalizers. OutcomeFinalizer handlers
// retain ownership of their separate result-aware completion lifecycle.
type InjectorFinalizer interface {
	Finalize(context.Context, InjectorLookup) error
}
