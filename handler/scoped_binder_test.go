package handler

import (
	"context"
	"testing"
)

type scopeABIBinder struct{}

func (scopeABIBinder) Bind(context.Context, any) error                     { return nil }
func (scopeABIBinder) Lookup(context.Context, ValueKey) (any, bool, error) { return nil, false, nil }
func (scopeABIBinder) BindScope(context.Context, any, string) error        { return nil }

type scopedABIContract struct{}

func (scopedABIContract) BindsInputByScope() bool { return true }
func TestOptionalScopedBindingPreservesExistingMethodValues(t *testing.T) {
	var b Binder = scopeABIBinder{}
	var s ScopedBinder = scopeABIBinder{}
	var bind func(context.Context, any) error = b.Bind
	var lookup func(context.Context, ValueKey) (any, bool, error) = b.Lookup
	var scoped func(context.Context, any, string) error = s.BindScope
	if bind(t.Context(), nil) != nil || scoped(t.Context(), nil, "") != nil {
		t.Fatal("method value")
	}
	_, _, _ = lookup(t.Context(), InputKey)
	var policy ScopedInputBinder = scopedABIContract{}
	if !policy.BindsInputByScope() {
		t.Fatal("policy")
	}
	var _ ContractFunc[int, int] = func(context.Context, Session, *int, *int) error { return nil }
}
