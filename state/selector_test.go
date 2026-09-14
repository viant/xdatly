package state

import "testing"

func TestSelectorAccessors(t *testing.T) {
	selector := &Selector{
		Limit:  25,
		Offset: 50,
		Page:   3,
	}

	if selector.CurrentLimit() != 25 {
		t.Fatalf("expected current limit")
	}
	if selector.CurrentOffset() != 50 {
		t.Fatalf("expected current offset")
	}
	if selector.CurrentPage() != 3 {
		t.Fatalf("expected current page")
	}
}

func TestSelectorsFind(t *testing.T) {
	selectors := Selectors{
		nil,
		&NamedSelector{Name: "vendor"},
		&NamedSelector{Name: "product"},
	}

	found := selectors.Find("product")
	if found == nil || found.Name != "product" {
		t.Fatalf("expected named selector lookup to succeed")
	}
}

func TestSelectorsClone(t *testing.T) {
	selectors := Selectors{&NamedSelector{
		Name: "products",
		Selector: Selector{
			Fields:       []string{"id"},
			Placeholders: []interface{}{1},
		},
	}}
	cloned := selectors.Clone()
	cloned[0].Fields[0] = "name"
	cloned[0].Placeholders[0] = 2
	if selectors[0].Fields[0] != "id" || selectors[0].Placeholders[0] != 1 {
		t.Fatalf("selector clone shares mutable slices: %+v", selectors[0])
	}
}
