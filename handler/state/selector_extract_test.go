package state

import "testing"

func TestExtractQuerySelector_Struct(t *testing.T) {
	type Spec struct {
		Page    int
		OrderBy string `json:"orderBy"`
		Limit   string `selector:"limit"`
	}

	qs, err := ExtractQuerySelector(Spec{Page: 2, OrderBy: "id desc", Limit: "10"})
	if err != nil {
		t.Fatalf("ExtractQuerySelector() error: %v", err)
	}
	if qs.Page != 2 {
		t.Fatalf("expected Page=2, got %d", qs.Page)
	}
	if qs.OrderBy != "id desc" {
		t.Fatalf("expected OrderBy=%q, got %q", "id desc", qs.OrderBy)
	}
	if qs.Limit != 10 {
		t.Fatalf("expected Limit=10, got %d", qs.Limit)
	}
}

func TestExtractQuerySelector_ParameterTag(t *testing.T) {
	type Spec struct {
		P int    `parameter:"page,kind=query,in=_page"`
		O string `parameter:"orderBy,kind=query,in=_orderby"`
	}

	qs, err := ExtractQuerySelector(Spec{P: 4, O: "id desc"})
	if err != nil {
		t.Fatalf("ExtractQuerySelector() error: %v", err)
	}
	if qs.Page != 4 {
		t.Fatalf("expected Page=4, got %d", qs.Page)
	}
	if qs.OrderBy != "id desc" {
		t.Fatalf("expected OrderBy=%q, got %q", "id desc", qs.OrderBy)
	}
}

func TestExtractNamedQuerySelector(t *testing.T) {
	named, err := ExtractNamedQuerySelector("v", struct{ Page int }{Page: 1})
	if err != nil {
		t.Fatalf("ExtractNamedQuerySelector() error: %v", err)
	}
	if named.Name != "v" {
		t.Fatalf("expected Name=v, got %q", named.Name)
	}
	if named.Page != 1 {
		t.Fatalf("expected Page=1, got %d", named.Page)
	}
}

func TestExtractQuerySelector_RejectsNonStruct(t *testing.T) {
	_, err := ExtractQuerySelector(map[string]any{"page": 1})
	if err == nil {
		t.Fatalf("expected error for map input")
	}
	_, err = ExtractQuerySelector(QuerySelector{Page: 1})
	if err == nil {
		t.Fatalf("expected error for QuerySelector input")
	}
}
