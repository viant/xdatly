package contract

import (
	"fmt"
	"go/token"
	"reflect"
	"strings"

	"github.com/viant/tagly/choice"
	"github.com/viant/tagly/tags"
)

// Field exposes only adjustable metadata, not type, layout or embedding flags.
// Helpers return copy-on-write metadata; callers must use the returned field.
// Helpers record errors for Validate so chained operations remain atomic from
// the host's perspective. Never apply a field for which Validate fails.
type Field struct {
	Name         string
	Tag          reflect.StructTag
	err          error
	skipChildren bool
	cloned       bool
}

// Clone creates an independent writable copy, even if f is already a copy.
// Use it before assigning Name/Tag directly or when branching from an edit.
func (f *Field) Clone() *Field { c := *f; c.cloned = true; return &c }

func (f *Field) writable() *Field {
	if f.cloned {
		return f
	}
	return f.Clone()
}
func (f *Field) withError(err error) *Field {
	if f.err != nil {
		return f
	}
	result := f.writable()
	result.err = err
	return result
}

// SkipChildren declines descendant callbacks for the policy returning this
// field. It changes no schema metadata and does not suppress nested policies.
func (f *Field) SkipChildren() *Field {
	if f.skipChildren {
		return f
	}
	result := f.writable()
	result.skipChildren = true
	return result
}
func (f *Field) ChildrenSkipped() bool { return f.skipChildren }

// Err reports a tag-edit error without inspecting a field's scalar type.
func (f *Field) Err() error { return f.err }

// Validate uses the original value type; enum choices never change that type
// or install runtime membership validation.
func (f *Field) Validate(originalType reflect.Type) error {
	if f.err != nil {
		return f.err
	}
	if !token.IsIdentifier(f.Name) || f.Name == "_" {
		return fmt.Errorf("invalid field name %q", f.Name)
	}
	if _, err := tags.Parse(string(f.Tag)); err != nil {
		return err
	}
	_, err := choice.Parse(f.Tag, originalType)
	return err
}

// SetTag replaces all occurrences of a key with one value, or adds a missing
// key. It preserves unrelated tags, including repeated and empty values.
func (f *Field) SetTag(key, value string) *Field { return f.set(key, []string{value}) }

// RemoveTag explicitly removes every occurrence. Removing an absent key is a no-op.
func (f *Field) RemoveTag(key string) *Field { return f.set(key, nil) }

func (f *Field) Hide() *Field {
	return f.SetTag("json", "-").SetTag("internal", "true").SetTag("mcp", "-")
}
func (f *Field) Required() *Field {
	return f.SetTag("required", "true").RemoveTag("optional").bindingRequired(true)
}
func (f *Field) Optional() *Field {
	return f.RemoveTag("required").RemoveTag("optional").bindingRequired(false)
}
func (f *Field) SetDescription(value string) *Field { return f.SetTag("description", value) }
func (f *Field) SetExample(value string) *Field     { return f.SetTag("example", value) }

// SetEnum replaces the complete ordered choice group. Values are textual tag
// literals; Validate interprets them using the original scalar type. Use
// RemoveTag("choice") to remove a domain, not SetEnum with an empty domain.
func (f *Field) SetEnum(values ...string) *Field {
	if len(values) == 0 {
		if f.err != nil {
			return f
		}
		return f.withError(fmt.Errorf("choice domain must not be empty"))
	}
	return f.set("choice", values)
}

func (f *Field) set(key string, values []string) *Field {
	if f.err != nil {
		return f
	}
	if !validTagKey(key) {
		return f.withError(fmt.Errorf("invalid tag key %q", key))
	}
	parsed, err := tags.Parse(string(f.Tag))
	if err != nil {
		return f.withError(err)
	}
	matched := 0
	equal := true
	for _, item := range parsed {
		if item.Name != key {
			continue
		}
		if matched >= len(values) || string(item.Values) != values[matched] {
			equal = false
		}
		matched++
	}
	if equal && matched == len(values) {
		return f
	}
	result := make(tags.Tags, 0, len(parsed)-matched+len(values))
	inserted := false
	appendValues := func() {
		for _, value := range values {
			result = append(result, &tags.Tag{Name: key, Values: tags.Values(value)})
		}
		inserted = true
	}
	for _, item := range parsed {
		if item.Name != key {
			result = append(result, item)
			continue
		}
		if !inserted {
			appendValues()
		}
	}
	if !inserted {
		appendValues()
	}
	updated := f.writable()
	updated.Tag = reflect.StructTag(result.Literal())
	return updated
}

func validTagKey(key string) bool {
	if key == "" {
		return false
	}
	for _, c := range key {
		if c <= ' ' || c == ':' || c == '"' || c == 0x7f {
			return false
		}
	}
	return true
}

// Keep the native binding annotation consistent with the public field metadata.
// The existing Tagly tokenizer preserves quoted/default values and their order.
func (f *Field) bindingRequired(required bool) *Field {
	if f.err != nil {
		return f
	}
	for _, key := range []string{"parameter", "bind"} {
		raw, ok := f.Tag.Lookup(key)
		if !ok {
			continue
		}
		name, rest := tags.Values(raw).Name()
		if strings.Contains(name, "=") {
			name, rest = "", tags.Values(raw)
		}
		var parts []string
		changed := false
		err := rest.MatchRawPairs(func(k, v string) error {
			if k == "required" {
				changed = true
				return nil
			}
			if v == "" {
				parts = append(parts, k)
			} else {
				parts = append(parts, k+"="+v)
			}
			return nil
		})
		if err != nil {
			return f.withError(err)
		}
		if required {
			parts = append(parts, "required=true")
			changed = true
		}
		if !changed {
			continue
		}
		value := strings.Join(parts, ",")
		if name != "" || strings.HasPrefix(raw, ",") {
			value = name + "," + value
		}
		f = f.SetTag(key, value)
	}
	return f
}
