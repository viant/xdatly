package bind

import (
	"context"
	"errors"
	"testing"

	xhandler "github.com/viant/xdatly/handler"
)

type testProvider struct {
	key   xhandler.ValueKey
	value any
	err   error
}

func (p testProvider) Key() xhandler.ValueKey {
	return p.key
}

func (p testProvider) Resolve(ctx context.Context) (any, error) {
	_ = ctx
	return p.value, p.err
}

type testBinder struct {
	values map[xhandler.ValueKey]any
	errs   map[xhandler.ValueKey]error
}

func (b *testBinder) Bind(_ context.Context, target any) error {
	_ = target
	return nil
}

func (b *testBinder) Lookup(_ context.Context, key xhandler.ValueKey) (any, bool, error) {
	if err, ok := b.errs[key]; ok {
		return nil, false, err
	}
	value, ok := b.values[key]
	return value, ok, nil
}

type testScope struct {
	binder xhandler.Binder
}

func (s testScope) Binder() xhandler.Binder {
	return s.binder
}

func TestProviderContractsCompile(t *testing.T) {
	var _ Provider = testProvider{}
	var _ Scope = testScope{}
}

func TestLookupTypedValue(t *testing.T) {
	binder := &testBinder{values: map[xhandler.ValueKey]any{xhandler.InputKey: "demo"}}
	value, ok, err := Lookup[string](context.Background(), binder, xhandler.InputKey)
	if err != nil {
		t.Fatalf("unexpected lookup error: %v", err)
	}
	if !ok || value != "demo" {
		t.Fatalf("expected typed lookup to return the stored string, got value=%q ok=%v", value, ok)
	}
}

func TestLookupUnknownKey(t *testing.T) {
	binder := &testBinder{values: map[xhandler.ValueKey]any{}}
	_, ok, err := Lookup[string](context.Background(), binder, xhandler.InputKey)
	if err != nil {
		t.Fatalf("unexpected unknown-key error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for unknown key")
	}
}

func TestLookupTypeMismatch(t *testing.T) {
	binder := &testBinder{values: map[xhandler.ValueKey]any{xhandler.InputKey: 42}}
	_, ok, err := Lookup[string](context.Background(), binder, xhandler.InputKey)
	if err == nil {
		t.Fatalf("expected type mismatch error")
	}
	if ok {
		t.Fatalf("expected ok=false on type mismatch")
	}
}

func TestLookupPropagatesBinderError(t *testing.T) {
	expected := errors.New("boom")
	binder := &testBinder{errs: map[xhandler.ValueKey]error{xhandler.InputKey: expected}}
	_, _, err := Lookup[string](context.Background(), binder, xhandler.InputKey)
	if !errors.Is(err, expected) {
		t.Fatalf("expected binder error propagation, got %v", err)
	}
}

func TestMustLookup(t *testing.T) {
	binder := &testBinder{values: map[xhandler.ValueKey]any{xhandler.InputKey: "demo"}}
	value, err := MustLookup[string](context.Background(), binder, xhandler.InputKey)
	if err != nil {
		t.Fatalf("unexpected must-lookup error: %v", err)
	}
	if value != "demo" {
		t.Fatalf("unexpected must-lookup value %q", value)
	}
}

func TestMustLookupUnknownKey(t *testing.T) {
	binder := &testBinder{values: map[xhandler.ValueKey]any{}}
	_, err := MustLookup[string](context.Background(), binder, xhandler.InputKey)
	if err == nil {
		t.Fatalf("expected unknown-key error")
	}
}
