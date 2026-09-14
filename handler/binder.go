package handler

import "context"

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
//
// Lookup returns (nil, false, nil) for an unknown key.
// An error is reserved for resolution failures of a known key.
type Binder interface {
	Bind(ctx context.Context, target any) error
	Lookup(ctx context.Context, key ValueKey) (any, bool, error)
}
