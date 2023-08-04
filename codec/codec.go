package codec

import (
	"context"
	"reflect"
)

type (
	Config struct {
		Body       string
		ParamType  reflect.Type
		Args       []string
		OutputType string
	}

	Factory interface {
		New(codecConfig *Config, options ...interface{}) (Instance, error)
	}

	Instance interface {
		ResultType(paramType reflect.Type) (reflect.Type, error)
		Value(ctx context.Context, raw interface{}, options ...interface{}) (interface{}, error)
	}
)
