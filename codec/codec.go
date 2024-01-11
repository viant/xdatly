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
		CriteriaBuilder func(ctx context.Context, expr string, input interface{}) (*Criteria, error) `json:"-" yaml:"-"`
	}

	Factory interface {
		New(codecConfig *Config, options ...Option) (Instance, error)
	}

	Instance interface {
		ResultType(inputType reflect.Type) (reflect.Type, error)
		Value(ctx context.Context, raw interface{}, options ...Option) (interface{}, error)
	}
)
