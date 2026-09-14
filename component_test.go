package xdatly

import "testing"

func TestComponentAuthoringPrimitiveCarriesTypedFields(t *testing.T) {
	type input struct {
		ID int
	}
	type output struct {
		Name string
	}

	component := Component[input, output]{
		Input:  input{ID: 7},
		Output: output{Name: "demo"},
	}

	if component.Input.ID != 7 {
		t.Fatalf("expected typed input field to be preserved")
	}
	if component.Output.Name != "demo" {
		t.Fatalf("expected typed output field to be preserved")
	}
}
