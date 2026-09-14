package predicate

import "testing"

func TestTemplateShape(t *testing.T) {
	template := &Template{
		Name:   "by_id",
		Source: "ID = ?",
		Args: []*NamedArgument{
			{Name: "id", Position: 0},
		},
	}

	if template.Name != "by_id" {
		t.Fatalf("expected template name to be preserved")
	}
	if template.Source != "ID = ?" {
		t.Fatalf("expected template source to be preserved")
	}
	if len(template.Args) != 1 || template.Args[0].Name != "id" || template.Args[0].Position != 0 {
		t.Fatalf("expected template args to be preserved")
	}
}

func TestNamedFilterShapes(t *testing.T) {
	filters := NamedFilters{
		&NamedFilter{Name: "country", Include: []string{"US"}, Exclude: []string{"CA"}},
	}
	if len(filters) != 1 {
		t.Fatalf("expected one named filter")
	}
	if filters[0].Name != "country" {
		t.Fatalf("expected filter name to be preserved")
	}
}
