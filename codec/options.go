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
	ColumnsSource
	Selector
	ValueGetter
	Options []interface{}
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

// WithColumnsSource creates columnSource option
func WithColumnsSource(columnSource ColumnsSource) Option {
	return func(o *Options) {
		o.ColumnsSource = columnSource
	}
}

// WithSelector creates selector option
func WithSelector(selector Selector) Option {
	return func(o *Options) {
		o.Selector = selector
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

// WithValueGetter creates value gettter options
func WithValueGetter(option ValueGetter) Option { //TODO replace with LookupValue
	return func(o *Options) {
		o.ValueGetter = option
	}
}

// WithOptions creates untyped options
func WithOptions(options ...interface{}) Option {
	return func(o *Options) {
		o.Options = options
	}
}
