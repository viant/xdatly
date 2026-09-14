package predicate

import (
	"context"
	"reflect"
	"testing"
)

func TestHandlerFunc(t *testing.T) {
	handler := HandlerFunc(func(_ context.Context, value any) (*Criteria, error) {
		return &Criteria{Expression: "id = ?", Placeholders: []any{value}}, nil
	})
	criteria, err := handler.Compute(context.Background(), 7)
	if err != nil {
		t.Fatalf("compute failed: %v", err)
	}
	if criteria.Expression != "id = ?" || !reflect.DeepEqual(criteria.Placeholders, []any{7}) {
		t.Fatalf("unexpected criteria: %#v", criteria)
	}
}
