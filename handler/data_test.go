package handler

import (
	"context"
	"reflect"
	"testing"
)

type testDataCapabilities struct{}

func (testDataCapabilities) Insert(string, any) error                            { return nil }
func (testDataCapabilities) Update(string, any) error                            { return nil }
func (testDataCapabilities) Delete(string, any) error                            { return nil }
func (testDataCapabilities) Execute(string, ...any) error                        { return nil }
func (testDataCapabilities) Allocate(context.Context, string, any, string) error { return nil }
func (testDataCapabilities) Flush(context.Context, string) error                 { return nil }

func TestDataCapabilitiesExposeFocusedAndCompleteSurfaces(t *testing.T) {
	var _ DML = testDataCapabilities{}
	var _ Sequencer = testDataCapabilities{}
	var _ Flusher = testDataCapabilities{}
	var _ Data = testDataCapabilities{}
	if DMLKey != "dml" || SequencerKey != "sequencer" || FlusherKey != "flusher" || DataKey != "data" {
		t.Fatalf("unexpected capability keys: %q %q %q %q", DMLKey, SequencerKey, FlusherKey, DataKey)
	}
	if methods := methodNames(reflect.TypeFor[DML]()); !reflect.DeepEqual(methods, []string{"Delete", "Execute", "Insert", "Update"}) {
		t.Fatalf("DML methods = %v", methods)
	}
	if methods := methodNames(reflect.TypeFor[Sequencer]()); !reflect.DeepEqual(methods, []string{"Allocate"}) {
		t.Fatalf("Sequencer methods = %v", methods)
	}
	if methods := methodNames(reflect.TypeFor[Flusher]()); !reflect.DeepEqual(methods, []string{"Flush"}) {
		t.Fatalf("Flusher methods = %v", methods)
	}
	if methods := methodNames(reflect.TypeFor[Data]()); !reflect.DeepEqual(methods, []string{"Allocate", "Delete", "Execute", "Flush", "Insert", "Update"}) {
		t.Fatalf("Data methods = %v", methods)
	}
}

func methodNames(aType reflect.Type) []string {
	result := make([]string, aType.NumMethod())
	for i := 0; i < aType.NumMethod(); i++ {
		result[i] = aType.Method(i).Name
	}
	return result
}
