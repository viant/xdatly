package codec

import (
	"context"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type testCodec struct{}

func (testCodec) Value(_ context.Context, raw interface{}, _ ...Option) (interface{}, error) {
	return raw, nil
}

type testFactory struct{}

func (testFactory) New(_ *Config, _ ...Option) (Instance, error) {
	return testCodec{}, nil
}

func TestCodecContractsCompile(t *testing.T) {
	var _ Instance = testCodec{}
	var _ Factory = testFactory{}
}

func TestOptionsHelpers(t *testing.T) {
	resourceMap := fstest.MapFS{"sql/query.sql": {Data: []byte("SELECT 1")}}
	resources := fs.FS(&resourceMap)
	options := NewOptions([]Option{
		WithRecord("demo"),
		WithOptions("x"),
		WithResourceFS(resources),
	})
	if options.Record != "demo" {
		t.Fatalf("expected record to be preserved")
	}
	if len(options.Options) != 1 || options.Options[0] != "x" {
		t.Fatalf("expected options payload to be preserved")
	}
	if options.ResourceFS != resources {
		t.Fatal("expected original resource filesystem to be preserved")
	}
}

func TestConfigSeparatesSourceAndDestinationTypes(t *testing.T) {
	type source struct{ ID int }
	type destination struct{ Values []int }
	config := Config{
		Body: "structql", SourceType: reflect.TypeOf([]source{}),
		DestinationType: reflect.TypeOf(destination{}), OutputTypeExpression: "destination",
	}
	if config.SourceType != reflect.TypeOf([]source{}) || config.DestinationType != reflect.TypeOf(destination{}) ||
		config.OutputTypeExpression != "destination" {
		t.Fatalf("codec config = %+v", config)
	}
}
