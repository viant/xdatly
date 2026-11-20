package handler

import (
	"context"
	"reflect"
)

// AssociationDelta is read-only transport intent for an opted-in collection.
// Members have the declared Go primitive type. Implementations must retain
// immutable, deduplicated sets; callers cannot mutate the captured request.
type AssociationDelta interface {
	AddedCount() int
	DeletedCount() int
	IsAdded(member any) bool
	IsDeleted(member any) bool
}

// AssociationDeltaLookup is supplied by the runtime for the exact invocation.
// Holder and path identify an individual relation occurrence, including its
// parent slice index. Nested component calls must not inherit another target's
// intent. This capability carries intent, not authorization or write permission.
type AssociationDeltaLookup interface {
	LookupAssociationDelta(target ComponentTarget, holder reflect.Type, field, path string) (AssociationDelta, bool)
}

const AssociationDeltaLookupKey ValueKey = "association_delta_lookup"

// LookupAssociationDelta returns trusted normalized intent, when present.
// Business hooks still own authorization, admission and dependent effects.
func LookupAssociationDelta(ctx context.Context, target ComponentTarget, holder reflect.Type, field, path string) (AssociationDelta, bool) {
	lookup, ok := ctx.Value(AssociationDeltaLookupKey).(AssociationDeltaLookup)
	if !ok || lookup == nil {
		return nil, false
	}
	return lookup.LookupAssociationDelta(target, holder, field, path)
}
