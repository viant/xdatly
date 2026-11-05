package contract

import (
	"reflect"
	"testing"

	"github.com/viant/tagly/choice"
)

func TestFieldUpsertAndExplicitRemoval(t *testing.T) {
	original := &Field{Name: "Mode", Tag: `sqlx:"MODE" json:"mode,omitempty" choice:"a" choice:"b" description:"old"`}
	changed := original.Clone().Hide().SetDescription("").SetExample(`a"b`).Required()
	if err := changed.Validate(reflect.TypeOf("")); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"json": "-", "internal": "true", "mcp": "-", "description": "", "example": `a"b`, "required": "true", "optional": "false", "sqlx": "MODE"} {
		if got, ok := changed.Tag.Lookup(key); !ok || got != want {
			t.Fatalf("%s=%q present=%v", key, got, ok)
		}
	}
	values, err := choice.Tags(changed.Tag)
	if err != nil || !reflect.DeepEqual(values, []string{"a", "b"}) {
		t.Fatal("sibling choices lost", values, err)
	}
	if original.Tag.Get("json") != "mode,omitempty" {
		t.Fatal("original mutated")
	}
	changed.Optional()
	if changed.Tag.Get("required") != "false" || changed.Tag.Get("optional") != "true" {
		t.Fatal("optional did not replace")
	}
	changed.RemoveTag("choice")
	if _, ok := changed.Tag.Lookup("choice"); ok {
		t.Fatal("choice occurrences remain")
	}
	before := changed.Tag
	changed.RemoveTag("absent").SetTag("json", "-")
	if changed.Tag != before {
		t.Fatal("no-op changed tag")
	}
}

func TestFieldEnumTypeOrderAndErrors(t *testing.T) {
	f := (&Field{Name: "Mode", Tag: `json:"mode" choice:"0" choice:"1"`}).SetEnum("9007199254740993", "3", "1")
	if err := f.Validate(reflect.TypeOf(int64(0))); err != nil {
		t.Fatal(err)
	}
	c, _ := choice.Parse(f.Tag, reflect.TypeOf(int64(0)))
	if !reflect.DeepEqual(c.WireValues(), []any{int64(9007199254740993), int64(3), int64(1)}) {
		t.Fatal(c)
	}
	for _, field := range []*Field{(&Field{Name: "X"}).SetEnum(), (&Field{Name: "X"}).SetEnum("1", "1"), (&Field{Name: "X"}).SetEnum("256"), (&Field{Name: "X"}).SetTag("bad:key", "value"), {Name: "X", Tag: `json:"unterminated`}} {
		if field.Validate(reflect.TypeOf(uint8(0))) == nil {
			t.Fatalf("accepted %+v", field)
		}
	}
	empty := (&Field{Name: "X"}).SetEnum("")
	if err := empty.Validate(reflect.TypeOf("")); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkFieldHide(b *testing.B) {
	f := Field{Name: "Value", Tag: `sqlx:"VALUE" json:"value,omitempty" description:"Value"`}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v := f.Clone().Hide()
		if v.Err() != nil {
			b.Fatal(v.Err())
		}
	}
}
