package predicate

import (
	"context"
	"reflect"
	"testing"
)

func TestVariadicHandlerAndPairs(t *testing.T) {
	var handler VariadicHandler = VariadicHandlerFunc(func(_ context.Context, value any, args ...any) (*Criteria, error) {
		options, err := Pairs(args...)
		if err != nil {
			return nil, err
		}
		return &Criteria{Placeholders: []any{value, options["name"], options["enabled"], options["limit"]}}, nil
	})
	got, err := handler.Compute(context.Background(), 7, "name", "p", "enabled", false, "limit", 10)
	if err != nil || !reflect.DeepEqual(got.Placeholders, []any{7, "p", false, 10}) {
		t.Fatalf("%#v %v", got, err)
	}
	if _, err = handler.Compute(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]any{{"odd"}, {1, "value"}, {"", "value"}, {"name", 1, "name", 2}} {
		if _, err := Pairs(args...); err == nil {
			t.Fatalf("accepted %#v", args)
		}
	}
}
