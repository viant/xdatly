package handler

import (
	"context"
	"errors"
	"testing"
)

type entityRecord struct{ Value int }
type entityParent struct{ ID int }
type entityPresence struct{ available, value bool }

func (p entityPresence) Has(field string) bool { return field == "Value" && p.value }
func (p entityPresence) Available() bool       { return p.available }

type entityHook struct{ initialized bool }

func (h *entityHook) Init(_ context.Context, value *entityRecord, state EntityState[entityRecord, entityParent]) error {
	if state.Parent == nil {
		return errors.New("parent required")
	}
	value.Value = state.Parent.ID
	h.initialized = true
	return nil
}
func (h *entityHook) Validate(_ context.Context, value *entityRecord, state EntityState[entityRecord, entityParent]) error {
	if !h.initialized || state.Previous == nil || !state.PreviousFields.Has("Value") {
		return errors.New("initialization and known previous value required")
	}
	if value.Value < state.Previous.Value {
		return errors.New("value decreased")
	}
	return nil
}

func TestEntityHooksShareTypedInvocationState(t *testing.T) {
	var hooks EntityHooks[entityRecord, entityParent] = &entityHook{}
	state := EntityState[entityRecord, entityParent]{Previous: &entityRecord{Value: 2}, Parent: &entityParent{ID: 3},
		Original: entityPresence{available: true}, PreviousFields: entityPresence{value: true}}
	value := &entityRecord{}
	if err := hooks.Init(context.Background(), value, state); err != nil {
		t.Fatal(err)
	}
	if err := hooks.Validate(context.Background(), value, state); err != nil {
		t.Fatal(err)
	}
	if value.Value != 3 || state.Previous.Value != 2 || !state.Original.Available() || state.Original.Has("Value") {
		t.Fatalf("entity=%+v state=%+v", value, state)
	}
	if err := (&entityHook{}).Validate(context.Background(), value, state); err == nil {
		t.Fatal("new hook instance lost phase state without failing")
	}
}

func TestOriginalPresenceDistinguishesMissingAndAbsent(t *testing.T) {
	for _, test := range []entityPresence{{}, {available: true}, {available: true, value: true}} {
		var presence OriginalPresence = test
		if presence.Available() != test.available || presence.Has("Value") != test.value || presence.Has("unknown") {
			t.Fatalf("presence=%+v", test)
		}
	}
}
