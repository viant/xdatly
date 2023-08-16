package codec

import "reflect"

// Options represents codec options
type Options struct {
	LookupType  func(name string) (reflect.Type, error)
	Record      interface{}
	LookupValue func(name string) (interface{}, bool)
	Options     []interface{}
}

func (o *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
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
func WithValueLookup(fn func(name string) (interface{}, bool)) Option {
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
