package contract

import (
	"reflect"
	"sync"
	"testing"
)

func TestFieldHelpersCopyOnceAndPreserveOriginal(t *testing.T) {
	original := &Field{Name: "Value", Tag: `json:"value" description:"Original"`}
	adjusted := original.Hide()
	if adjusted == original || !adjusted.cloned || original.cloned {
		t.Fatal("first write did not create a writable copy")
	}
	if original.Tag != `json:"value" description:"Original"` || original.Err() != nil || original.ChildrenSkipped() {
		t.Fatal("source modified")
	}
	chained := adjusted.SetDescription("Hidden").SetExample("sample").Required().Optional().SetEnum("a", "b").RemoveTag("choice").SkipChildren()
	if chained != adjusted {
		t.Fatal("a chained helper cloned an already writable field")
	}
	if chained.Tag.Get("json") != "-" || chained.Tag.Get("description") != "Hidden" || !chained.ChildrenSkipped() {
		t.Fatalf("%+v", chained)
	}
	if err := chained.Validate(reflect.TypeFor[string]()); err != nil {
		t.Fatal(err)
	}
	if original.Tag.Get("json") != "value" || original.ChildrenSkipped() {
		t.Fatal("chain changed original")
	}
}

func TestFieldNoOpHelpersRetainIdentity(t *testing.T) {
	field := &Field{Name: "Value", Tag: `json:"value" description:"Same" choice:"a" choice:"b"`}
	for _, apply := range []func(*Field) *Field{
		func(f *Field) *Field { return f.SetTag("json", "value") },
		func(f *Field) *Field { return f.SetDescription("Same") },
		func(f *Field) *Field { return f.RemoveTag("missing") },
		func(f *Field) *Field { return f.Optional() },
		func(f *Field) *Field { return f.SetEnum("a", "b") },
	} {
		if actual := apply(field); actual != field || actual.cloned {
			t.Fatal("no-op copied metadata")
		}
	}
	hidden := &Field{Name: "Value", Tag: `json:"-" internal:"true" mcp:"-"`}
	if hidden.Hide() != hidden {
		t.Fatal("already-hidden field cloned")
	}
	skipped := field.SkipChildren()
	if skipped.SkipChildren() != skipped {
		t.Fatal("skip flag cloned twice")
	}
}

func TestFieldBindingRequirementCanBeFirstWrite(t *testing.T) {
	original := &Field{Name: "Value", Tag: `json:"value" parameter:",kind=query,in=value,required=true,value='a,b'" bind:",kind=query,in=value,required=true"`}
	optional := original.Optional()
	if optional == original || optional.Tag.Get("parameter") != ",kind=query,in=value,value='a,b'" || optional.Tag.Get("bind") != ",kind=query,in=value" {
		t.Fatalf("%s", optional.Tag)
	}
	if original.Tag.Get("parameter") != ",kind=query,in=value,required=true,value='a,b'" {
		t.Fatal("binding source mutated")
	}
	requiredSource := &Field{Name: "Value", Tag: `required:"true" parameter:",kind=query,in=value" bind:",kind=query,in=value"`}
	required := requiredSource.Required()
	if required == requiredSource || required.Tag.Get("parameter") != ",kind=query,in=value,required=true" || required.Tag.Get("bind") != ",kind=query,in=value,required=true" {
		t.Fatalf("%s", required.Tag)
	}
}

func TestFieldCopyOnWriteErrorIsolationAndBranching(t *testing.T) {
	original := &Field{Name: "Value", Tag: `json:"value"`}
	bad := original.SetTag("bad:key", "value")
	if bad == original || bad.Err() == nil || original.Err() != nil || bad.Tag != original.Tag {
		t.Fatal("helper error contaminated source")
	}
	if bad.SetExample("ignored") != bad {
		t.Fatal("errored chain cloned again")
	}
	empty := original.SetEnum()
	if empty == original || empty.Err() == nil || original.Err() != nil {
		t.Fatal("enum error contaminated source")
	}
	first := original.Clone()
	if !first.cloned || first.SetDescription("first") != first {
		t.Fatal("explicit Clone not writable")
	}
	second := first.Clone()
	second.Name = "Renamed"
	second = second.SetDescription("second")
	if second == first || first.Name != "Value" || first.Tag.Get("description") != "first" {
		t.Fatal("explicit branch mutated earlier edit")
	}
}

func TestConcurrentFieldEditsStartFromReadOnlySource(t *testing.T) {
	original := &Field{Name: "Value", Tag: `json:"value"`}
	var group sync.WaitGroup
	for i := 0; i < 16; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			field := original.Hide().Required()
			if field == original || field.Err() != nil || field.Tag.Get("required") != "true" {
				t.Error("invalid independent edit")
			}
		}()
	}
	group.Wait()
	if original.Tag != `json:"value"` || original.cloned {
		t.Fatal("read-only baseline mutated")
	}
}

func BenchmarkFieldCopyOnWrite(b *testing.B) {
	for _, explicit := range []bool{false, true} {
		name := "automatic"
		if explicit {
			name = "explicit"
		}
		b.Run(name, func(b *testing.B) {
			original := &Field{Name: "Value", Tag: `json:"value"`}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				field := original
				if explicit {
					field = field.Clone()
				}
				field = field.Hide().SetDescription("Hidden").Required()
				if field.Err() != nil {
					b.Fatal(field.Err())
				}
			}
		})
	}
}
