// Package contract defines caller-dependent adjustments to a component's public
// contract. Binding, traversal, caching and codec execution belong to the host.
package contract

import (
	"context"
	"reflect"

	"github.com/viant/bindly"
)

// Provider may be implemented by an original linked view/domain type,
// component Input/Output, or composite. Nested providers are discovered by the host.
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
	// DiscoveryAlternatives declares complete finite output policies in place of
	// this owner's selected policy for discovery only. Each alternative retains
	// canonical storage and composes with every other nested/enclosing owner.
	// Construct alternatives from the same authorization decision as Field; the
	// host does not invoke Contract again for each branch. Calls use only Field.
	// Nil/empty retains the selected policy. Nil entries and recursive alternatives
	// are invalid. Names and NoCache follow the ordinary variant rules.
	DiscoveryAlternatives []*Policy
	Name                  string
	// NoCache opts out of host variant reuse even when Name is set. It does not
	// disable unrelated codec/type caches or cache authorization decisions.
	NoCache bool
	// Field receives the original containing struct type and its canonical path
	// (excluding the field name). Use helper return values for copy-on-write edits; use Clone before assigning
	// Name or Tag directly.
	// Nil or an equivalent replacement retains the original field. Only Name
	// and Tag can change. The host must check replacement.Validate(originalType)
	// before publishing or caching the effective contract.
	Field func(parentPath string, parentType reflect.Type, field *Field) (*Field, error)
}
