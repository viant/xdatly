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

type entityOutput struct{ Violations []string }
type entityHook struct{ initialized bool }

func (h *entityHook) Init(_ context.Context, value *entityRecord, state LifecycleContext[entityRecord, entityParent, entityOutput]) error {
	if state.Parent == nil || state.Output == nil {
		return errors.New("parent and output required")
	}
	value.Value = state.Parent.ID
	h.initialized = true
	return nil
}
func (h *entityHook) Validate(_ context.Context, value *entityRecord, state LifecycleContext[entityRecord, entityParent, entityOutput]) error {
	if !h.initialized || state.Previous == nil || !state.PreviousFields.Has("Value") {
		return errors.New("initialization and known previous value required")
	}
	if value.Value < state.Previous.Value {
		return errors.New("value decreased")
	}
	state.Output.Violations = append(state.Output.Violations, "observed")
	return nil
}

func TestEntityHooksShareTypedInvocationStateAndOutput(t *testing.T) {
	var hooks EntityHooks[entityRecord, entityParent, entityOutput] = &entityHook{}
	output := &entityOutput{}
	state := LifecycleContext[entityRecord, entityParent, entityOutput]{
		EntityState: EntityState[entityRecord, entityParent]{Previous: &entityRecord{Value: 2}, Parent: &entityParent{ID: 3},
			Original: entityPresence{available: true}, PreviousFields: entityPresence{value: true}},
		Output: output,
	}
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
	if len(output.Violations) != 1 || output.Violations[0] != "observed" {
		t.Fatalf("output=%+v", output)
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

type lifecycleContextEntity struct{ Name string }
type lifecycleContextParent struct{ ID int }
type lifecycleContextOutput struct{ Violations []string }

func TestLifecycleContextEmbedsEntityStateAndSharesOutput(t *testing.T) {
	previous := &lifecycleContextEntity{Name: "before"}
	parent := &lifecycleContextParent{ID: 7}
	output := &lifecycleContextOutput{}
	state := LifecycleContext[lifecycleContextEntity, lifecycleContextParent, lifecycleContextOutput]{
		EntityState: EntityState[lifecycleContextEntity, lifecycleContextParent]{
			Previous: previous,
			Parent:   parent,
		},
		Output: output,
	}

	if state.Previous != previous || state.Parent != parent {
		t.Fatalf("embedded entity state was not preserved: %+v", state)
	}
	state.Output.Violations = append(state.Output.Violations, "invalid name")
	if len(output.Violations) != 1 || output.Violations[0] != "invalid name" {
		t.Fatalf("lifecycle output was not shared: %+v", output)
	}
}
