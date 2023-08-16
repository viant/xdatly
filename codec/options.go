package codec

import (
	"context"
	"reflect"
)

// Options represents codec options
type Options struct {
	LookupType  func(name string) (reflect.Type, error)
	Record      interface{}
	LookupValue func(ctx context.Context, name string) (interface{}, error)
	Options     []interface{}
}

func (o *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// NewOptions creates an options
func NewOptions(opts []Option) *Options {
	options := &Options{}
	options.Apply(opts)
	return options
}

// Option func to set options
type Option func(o *Options)

// WithTypeLookup creates type lookup option
func WithTypeLookup(fn func(name string) (reflect.Type, error)) Option {
	return func(o *Options) {
		o.LookupType = fn
	}
}

// WithRecord creates type record option
func WithRecord(record interface{}) Option {
	return func(o *Options) {
		o.Record = record
	}
}

// WithValueLookup creates value lookup option
func WithValueLookup(fn func(ctx context.Context, name string) (interface{}, error)) Option {
	return func(o *Options) {
		o.LookupValue = fn
	}
}

// WithOptions creates untyped options
func WithOptions(options ...interface{}) Option {
	return func(o *Options) {
		o.Options = options
	}
}
