package codec

import (
	"context"
	"reflect"
)

type (
	Config struct {
		Body       string
		InputType  reflect.Type
		Args       []string
		OutputType string
	}

	Factory interface {
		New(codecConfig *Config, options ...Option) (Instance, error)
	}

	Instance interface {
		ResultType(inputType reflect.Type) (reflect.Type, error)
		Value(ctx context.Context, raw interface{}, options ...Option) (interface{}, error)
	}
)
