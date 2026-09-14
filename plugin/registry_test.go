package plugin

import (
	"context"
	"testing"

	xcodec "github.com/viant/xdatly/codec"
	xdocs "github.com/viant/xdatly/docs"
	xpredicate "github.com/viant/xdatly/predicate"
)

type testCodec struct{}

func (testCodec) Value(_ context.Context, raw interface{}, _ ...xcodec.Option) (interface{}, error) {
	return raw, nil
}

type testFactory struct{}

func (testFactory) New(_ *xcodec.Config, _ ...xcodec.Option) (xcodec.Instance, error) {
	return testCodec{}, nil
}

type testDocsProvider struct{}

func (testDocsProvider) Service(_ context.Context, _ ...xdocs.Option) (xdocs.Service, error) {
	return nil, nil
}

type testRegistry struct {
	codecName     string
	factoryName   string
	predicateName string
	docsName      string
}

func (r *testRegistry) RegisterCodec(name string, _ xcodec.Instance)       { r.codecName = name }
func (r *testRegistry) RegisterCodecFactory(name string, _ xcodec.Factory) { r.factoryName = name }
func (r *testRegistry) RegisterPredicate(template *xpredicate.Template) {
	if template != nil {
		r.predicateName = template.Name
	}
}
func (r *testRegistry) RegisterDocs(name string, _ xdocs.Provider) { r.docsName = name }

func TestPluginRegistryContract(t *testing.T) {
	var _ Registry = &testRegistry{}

	reg := &testRegistry{}
	reg.RegisterCodec("jwt", testCodec{})
	reg.RegisterCodecFactory("json", testFactory{})
	reg.RegisterPredicate(&xpredicate.Template{Name: "by_id"})
	reg.RegisterDocs("default", testDocsProvider{})

	if reg.codecName != "jwt" || reg.factoryName != "json" || reg.predicateName != "by_id" || reg.docsName != "default" {
		t.Fatalf("expected registry methods to preserve names")
	}
}
