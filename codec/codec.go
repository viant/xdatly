package codec

import (
	"context"
	"reflect"
)

type (
	Config struct {
		Body       string
		InputType  reflect.Type `json:"-" yaml:"-"`
		Args       []string
		OutputType string
		//Optional builder
	}

	Factory interface {
		New(codecConfig *Config, options ...Option) (Instance, error)
	}

	Instance interface {
		ResultType(inputType reflect.Type) (reflect.Type, error)
		Value(ctx context.Context, raw interface{}, options ...Option) (interface{}, error)
	}
)
