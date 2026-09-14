package bind

import (
	"context"
	"fmt"

	xhandler "github.com/viant/xdatly/handler"
)

// Provider resolves one binder-scoped value. Implementations remain runtime-
// owned; this package only defines the public contract they target.
type Provider interface {
	Key() xhandler.ValueKey
	Resolve(ctx context.Context) (any, error)
}

// Scope exposes the binder-backed resolution surface for one invocation or
// other bounded value scope.
type Scope interface {
	Binder() xhandler.Binder
}

// Note: this package intentionally does not expose a generic reflection-based
// injection helper yet. Current Datly-side injection is still runtime-colored
// and type-specific (logger / validator / sequencer / message bus wiring), so
// promoting that logic here would overfit one runtime rather than define a
// stable cross-runtime contract.

// Lookup resolves a binder value and type-asserts it to T.
// Unknown keys return (zero, false, nil).
func Lookup[T any](ctx context.Context, binder xhandler.Binder, key xhandler.ValueKey) (T, bool, error) {
	var zero T
	if binder == nil {
		return zero, false, nil
	}
	value, ok, err := binder.Lookup(ctx, key)
	if err != nil || !ok {
		return zero, ok, err
	}
	typed, cast := value.(T)
	if !cast {
		return zero, false, fmt.Errorf("binder value %q has type %T, not %T", key, value, zero)
	}
	return typed, true, nil
}

// MustLookup resolves a binder value, returning an error when the key is
// unknown or when the resolved value does not match T.
func MustLookup[T any](ctx context.Context, binder xhandler.Binder, key xhandler.ValueKey) (T, error) {
	var zero T
	value, ok, err := Lookup[T](ctx, binder, key)
	if err != nil {
		return zero, err
	}
	if !ok {
		return zero, fmt.Errorf("binder value %q is not available", key)
	}
	return value, nil
}
