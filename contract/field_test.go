package contract

import (
	"reflect"
	"testing"

	"github.com/viant/bindly"
	"github.com/viant/tagly/choice"
)

func TestOptionalRemovesNativeBindingRequirementWithoutLosingMetadata(t *testing.T) {
	for _, key := range []string{"parameter", "bind"} {
		f := (&Field{Name: "Value"}).SetTag(key, `,kind=query,in=value,required=true,value='a,b'`).SetTag("required", "true")
		f = f.Optional()
		if _, ok := f.Tag.Lookup("required"); ok {
			t.Fatal("explicit flag remains")
		}
		binding, ok, err := bindly.BindingSpecFromField(reflect.StructField{Name: f.Name, Type: reflect.TypeOf(""), Tag: f.Tag})
		if err != nil || !ok {
			t.Fatal(err)
		}
		if binding.Required != nil && *binding.Required {
			t.Fatal("binding still required")
		}
		if binding.Location.In != "value" || binding.DefaultValue != "a,b" {
			t.Fatalf("binding metadata changed: %+v", binding)
		}
		f = f.Required()
		binding, _, err = bindly.BindingSpecFromField(reflect.StructField{Name: f.Name, Type: reflect.TypeOf(""), Tag: f.Tag})
		if err != nil || binding.Required == nil || !*binding.Required {
			t.Fatal("required binding not restored", err)
		}
	}
}

func TestFieldUpsertAndExplicitRemoval(t *testing.T) {
	original := &Field{Name: "Mode", Tag: `sqlx:"MODE" json:"mode,omitempty" choice:"a" choice:"b" description:"old"`}
	changed := original.Clone().Hide().SetDescription("").SetExample(`a"b`).Required()
	if err := changed.Validate(reflect.TypeOf("")); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"json": "-", "internal": "true", "mcp": "-", "description": "", "example": `a"b`, "required": "true", "sqlx": "MODE"} {
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
	changed = changed.Optional()
	if _, found := changed.Tag.Lookup("required"); found {
		t.Fatal("optional retained explicit requiredness")
	}
	if _, found := changed.Tag.Lookup("optional"); found {
		t.Fatal("optional default should need no flag")
	}
	changed = changed.RemoveTag("choice")
	if _, ok := changed.Tag.Lookup("choice"); ok {
		t.Fatal("choice occurrences remain")
	}
	before := changed.Tag
	changed = changed.RemoveTag("absent").SetTag("json", "-")
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

func TestOptionalFirstBindingOption(t *testing.T) {
	for _, tag := range []reflect.StructTag{
		`bind:"required=true,kind=query,in=value"`,
		`parameter:"Value,required=true,kind=query,in=value"`,
		`bind:"kind=query,in=value,required=true"`,
	} {
		f := (&Field{Name: "Value", Tag: tag}).Optional()
		binding, ok, err := bindly.BindingSpecFromField(reflect.StructField{Name: f.Name, Type: reflect.TypeFor[string](), Tag: f.Tag})
		if err != nil || !ok {
			t.Fatalf("%s: %v", f.Tag, err)
		}
		if binding.Required != nil && *binding.Required || binding.Location.Kind != "query" || binding.Location.In != "value" {
			t.Fatalf("lost binding metadata: %+v", binding)
		}
	}
}

func BenchmarkFieldOptional(b *testing.B) {
	for _, tc := range []struct {
		name string
		tag  reflect.StructTag
	}{
		{"unchanged", `json:"value"`},
		{"binding", `json:"value" required:"true" parameter:",kind=query,in=value,required=true,value='a,b'"`},
	} {
		b.Run(tc.name, func(b *testing.B) {
			f := Field{Name: "Value", Tag: tc.tag}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if err := f.Clone().Optional().Err(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func TestSkipChildrenIsSeparateFromFieldMetadata(t *testing.T) {
	original := &Field{Name: "Child", Tag: `json:"child"`}
	selected := original.Clone().SkipChildren()
	if !selected.ChildrenSkipped() || original.ChildrenSkipped() || selected.Name != original.Name || selected.Tag != original.Tag {
		t.Fatal("skip mutated schema metadata or original field")
	}
	if !selected.Clone().Hide().ChildrenSkipped() {
		t.Fatal("chained edit lost traversal decision")
	}
	if err := selected.Validate(reflect.TypeFor[string]()); err != nil {
		t.Fatal(err)
	}
}
