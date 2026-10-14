package response_test

import (
	"encoding/json"
	legacy "github.com/viant/xdatly/handler/response"
	modern "github.com/viant/xdatly/response"
	"reflect"
	"testing"
)

func requireComparable[T comparable](v T) T { return v }
func TestOriginalStatusNominalContractAndModernSeparation(t *testing.T) {
	got := reflect.TypeFor[legacy.Status]()
	if got.PkgPath() != "github.com/viant/xdatly/handler/response" || got.Name() != "Status" || got.NumField() != 5 || got.NumMethod() != 0 {
		t.Fatalf("wrong original authority: %v", got)
	}
	names := []string{"Status", "Message", "Errors", "Warning", "Extras"}
	types := []reflect.Type{reflect.TypeFor[string](), reflect.TypeFor[string](), reflect.TypeFor[interface{}](), reflect.TypeFor[[]*legacy.Warning](), reflect.TypeFor[map[string]interface{}]()}
	tags := []string{`json:"status,omitempty"`, `json:"message,omitempty"`, `json:"errors,omitempty"`, `json:"warning,omitempty"`, `json:"extras,omitempty" format:"inline=true"`}
	for i, name := range names {
		f := got.Field(i)
		if f.Name != name || f.Type != types[i] || string(f.Tag) != tags[i] {
			t.Fatalf("original field%d changed: %+v", i, f)
		}
	}
	w := reflect.TypeFor[legacy.Warning]()
	if w.NumField() != 2 || w.PkgPath() != got.PkgPath() || w.NumMethod() != 0 {
		t.Fatal(w)
	}
	if w.Field(0).Name != "Message" || string(w.Field(0).Tag) != `json:"message,omitempty"` || w.Field(1).Name != "Reason" || string(w.Field(1).Tag) != `json:"reason,omitempty"` {
		t.Fatal("Warning source fields changed")
	}
	original := legacy.Status{"ok", "source", false, []*legacy.Warning{{"warning", "reason"}}, map[string]interface{}{"zero": 0}}
	type envelope struct {
		legacy.Status
		Data []int `json:"data"`
	}
	out := envelope{Status: original, Data: []int{7}}
	if out.Errors != false || out.Warning[0].Reason != "reason" || out.Extras["zero"] != 0 || out.Status.Status != "ok" {
		t.Fatal("original promoted selectors lost")
	}
	clone := original
	clone.Warning = append([]*legacy.Warning(nil), original.Warning...)
	clone.Extras = map[string]interface{}{}
	for k, v := range original.Extras {
		clone.Extras[k] = v
	}
	if !reflect.DeepEqual(clone, original) {
		t.Fatal("source copy use changed")
	}
	m := modern.Status{"ok", "modern", 42}
	if requireComparable(m) != m || m != (modern.Status{"ok", "modern", 42}) {
		t.Fatal("modern comparable contract changed")
	}
	keys := map[modern.Status]bool{m: true}
	if !keys[m] {
		t.Fatal("modern map key changed")
	}
	mt := reflect.TypeFor[modern.Status]()
	if mt.PkgPath() != "github.com/viant/xdatly/response" || mt.NumField() != 3 || mt == got || mt.AssignableTo(got) || got.AssignableTo(mt) {
		t.Fatal("nominal modern/legacy distinction lost")
	}
	for i, n := range []string{"Status", "Message", "Error"} {
		if mt.Field(i).Name != n {
			t.Fatal("modern fields changed")
		}
	}
	raw, e := json.Marshal(out)
	if e != nil {
		t.Fatal(e)
	}
	var wire map[string]json.RawMessage
	if e = json.Unmarshal(raw, &wire); e != nil {
		t.Fatal(e)
	}
	if string(wire["errors"]) != "false" || string(wire["extras"]) != `{"zero":0}` || string(wire["data"]) != "[7]" {
		t.Fatalf("original stdlib semantics changed: %s", raw)
	}
}
