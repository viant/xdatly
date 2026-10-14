package response

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestOriginalNominalErrorDataAndMethodContract(t *testing.T) {
	typ := reflect.TypeFor[Error]()
	if typ.PkgPath() != "github.com/viant/xdatly/handler/response" || typ.Name() != "Error" || typ.NumField() != 6 {
		t.Fatal("nominal identity/field count", typ)
	}
	fs := []struct {
		Name string
		Type reflect.Type
		Tag  reflect.StructTag
	}{{"View", reflect.TypeFor[string](), `json:"view,omitempty" `}, {"Parameter", reflect.TypeFor[string](), `json:"parameter,omitempty" `}, {"Code", reflect.TypeFor[int](), `json:"statusCode,omitempty" `}, {"Err", reflect.TypeFor[error](), `json:"-"`}, {"Message", reflect.TypeFor[string](), `json:"message,omitempty" `}, {"Object", reflect.TypeFor[any](), `json:"object,omitempty" `}}
	for i, w := range fs {
		f := typ.Field(i)
		if f.Name != w.Name || f.Type != w.Type || f.Tag != w.Tag {
			t.Fatalf("field%d %+v want%+v", i, f, w)
		}
	}
	ptr := reflect.PointerTo(typ)
	if ptr.NumMethod() != 4 {
		t.Fatal("expanded method set", ptr.NumMethod())
	}
	for _, name := range []string{"StatusCode", "SetStatusCode", "Reader", "Error"} {
		if _, ok := ptr.MethodByName(name); !ok {
			t.Fatal("missing method", name)
		}
	}
	var contracts interface {
		error
		Reader() io.Reader
		StatusCode() int
		SetStatusCode(int)
	} = &Error{}
	if _, ok := contracts.(interface{ Unwrap() error }); ok {
		t.Fatal("legacy error gained wrapping behavior")
	}
	if _, ok := contracts.(interface{ ResponseBody() any }); ok {
		t.Fatal("legacy error gained explicit body behavior")
	}
}
func TestOriginalNominalErrorPublicMessageAndHiddenCause(t *testing.T) {
	cause := errors.New("private SQL and credential diagnostic")
	e := &Error{Parameter: "records", Code: 401, Message: "public", Err: cause, Object: []*Error{{Parameter: "tenant", Code: 401, Message: "nested", Err: cause}}}
	if e.Error() != "public" || e.StatusCode() != 401 || e.Err != cause {
		t.Fatal("original behavior changed")
	}
	raw, err := io.ReadAll(e.Reader())
	if err != nil || string(raw) != "public" {
		t.Fatal("reader", err)
	}
	wire, err := json.Marshal(e)
	if err != nil || strings.Contains(string(wire), "private") || strings.Contains(string(wire), "view") {
		t.Fatal("hidden cause/omission", string(wire), err)
	}
	e.SetStatusCode(403)
	if e.Code != 403 {
		t.Fatal("status setter")
	}
	zero, err := json.Marshal(&Error{Err: cause})
	if err != nil || string(zero) != "{}" {
		t.Fatal("zero contract", string(zero), err)
	}
}
