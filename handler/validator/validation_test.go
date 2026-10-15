package validator_test

import (
	"encoding/json"
	modern "github.com/viant/xdatly/handler"
	legacy "github.com/viant/xdatly/handler/validator"
	"reflect"
	"testing"
)

var keyed = legacy.Violation{Location: "Data[0].Id", Field: "Id", Value: 0, Message: "missing", Check: "required"}
var unkeyed = legacy.Violation{"Data[0].Id", "Id", 0, "missing", "required"}
var named legacy.Violations = []*legacy.Violation{&keyed, &unkeyed}
var _ []*legacy.Violation = named

func TestOriginalNominalFieldsAndTags(t *testing.T) {
	typ := reflect.TypeFor[legacy.Violation]()
	if typ.Name() != "Violation" || typ.PkgPath() != "github.com/viant/xdatly/handler/validator" {
		t.Fatal(typ)
	}
	names := []string{"Location", "Field", "Value", "Message", "Check"}
	for i, n := range names {
		f := typ.Field(i)
		if f.Name != n || string(f.Tag) != `json:",omitempty"` {
			t.Fatalf("field %d %v", i, f)
		}
	}
	if typ.NumField() != 5 || typ.NumMethod() != 0 {
		t.Fatal(typ)
	}
	if !reflect.DeepEqual(keyed, unkeyed) {
		t.Fatal("literal contract differs")
	}
}
func TestOriginalPureMethods(t *testing.T) {
	v := legacy.NewValidation()
	if v.Violations == nil || len(v.Violations) != 0 || v.Failed {
		t.Fatal(v)
	}
	v.Append("Data[2].Id", "Id", 23, "required", "")
	if !v.Failed || v.Violations[0].Message != "check 'required' failed on field Id" {
		t.Fatal(v)
	}
	v.AddViolation("Data[1].Id", "Id", 17, "custom", "ID $value / $value")
	if v.Violations[1].Message != "ID 17 / $value" {
		t.Fatal(v.Violations[1])
	}
	before := v.Violations
	v.Violations.Sort()
	if &before[0] != &v.Violations[0] || v.Violations[0].Location != "Data[1].Id" || v.Violations[1].Location != "Data[2].Id" {
		t.Fatal(v)
	}
}
func TestOriginalInterfaceOmissionSemantics(t *testing.T) {
	var p *int
	for _, tc := range []struct {
		name  string
		value any
		want  string
	}{
		{"nil", nil, `{}`}, {"zero", 0, `{"Value":0}`}, {"false", false, `{"Value":false}`}, {"empty-string", "", `{"Value":""}`}, {"typed-nil", p, `{"Value":null}`}, {"nil-slice", []int(nil), `{"Value":null}`}, {"empty-slice", []int{}, `{"Value":[]}`}, {"value", 3, `{"Value":3}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(legacy.Violation{Value: tc.value})
			if err != nil || string(raw) != tc.want {
				t.Fatalf("%s %v want %s", raw, err, tc.want)
			}
		})
	}
}
func TestModernValidationUnchanged(t *testing.T) {
	v := &modern.Validation{Violations: []*modern.Violation{{Location: "Input.Id", Field: "Id", CheckedValue: "private-sentinel", HasCheckedValue: true}}}
	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err = json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"violations":[{"location":"Input.Id","field":"Id"}]}` || v.StatusCode() != 422 || v.Err() != v {
		t.Fatalf("%s %d", raw, v.StatusCode())
	}
}
