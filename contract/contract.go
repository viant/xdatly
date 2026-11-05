// Package contract defines caller-dependent adjustments to a component's public
// contract. Binding, traversal, caching and codec execution belong to the host.
package contract

import (
	"context"
	"reflect"

	"github.com/viant/bindly"
)

// Provider may be implemented by a component's original linked Input or Output.
// The binder belongs to the current invocation. A provider must not retain it.
// Return nil, nil to retain the original linked contract without traversal.
type Provider interface {
	Contract(context.Context, bindly.Binder) (*Policy, error)
}

// Policy describes an immutable adjustment decision. Name identifies equivalent
// decisions within the host's component/type/direction/generation namespace.
// An empty Name disables retained variant caching. A nil Field makes no changes.
// Hosts cache effective types/plans, never this policy's captured request state.
type Policy struct {
	Name string
	// Field receives the original containing struct type and its canonical path
	// (excluding the field name). The supplied field is read-only; use Clone.
	// Nil or an equivalent replacement retains the original field. Only Name
	// and Tag can change. The host must check replacement.Validate(originalType)
	// before publishing or caching the effective contract.
	Field func(parentPath string, parentType reflect.Type, field *Field) (*Field, error)
}
