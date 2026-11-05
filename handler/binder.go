package handler

import (
	"context"

	"github.com/viant/bindly"
)

type ValueKey string

const (
	// InputKey is the reserved canonical component input key exposed to hooks,
	// predicates, and compatibility paths that need direct input visibility.
	InputKey ValueKey = "input"
	// SelectorsKey is the invocation-scoped dynamic view-selector capability.
	// Values are resolved through the Binder; readers must not carry a parallel
	// selector input channel.
	SelectorsKey ValueKey = "selectors"
)

// Binder resolves and binds values for one invocation scope.
// It is intentionally narrow and does not expose service construction.
// Bind is the shared Bindly interface. Existing two-argument calls are unchanged;
// specialized invocation binders may reject options that their scope cannot honor.
//
// Lookup returns (nil, false, nil) for an unknown key.
// An error is reserved for resolution failures of a known key.
type Binder interface {
	bindly.Binder
	Lookup(ctx context.Context, key ValueKey) (any, bool, error)
}

// ScopedInputBinder is a static opt-in on a typed custom contract. The engine
// binds its exact default scope before Exec; named scopes remain handler-selected.
type ScopedInputBinder interface{ BindsInputByScope() bool }

// ScopedBinder selects an exact declared scope on the canonical input. Bind
// selects the empty scope again. It is valid only during the handler invocation.
type ScopedBinder interface {
	Binder
	BindScope(ctx context.Context, target any, scope string) error
}
