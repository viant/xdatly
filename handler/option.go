package handler

import (
	"reflect"
)

type (
	Options struct {
		Arguments  []string
		InputType  reflect.Type
		OutputType reflect.Type
		LookupType func(name string) (reflect.Type, error)
	}

	Option func(o *Options)
)

// WithArguments creates arguments option
func WithArguments(args []string) Option {
	return func(o *Options) {
		o.Arguments = args
	}
}

// WithInputType creates inputType option
func WithInputType(inputType reflect.Type) Option {
	return func(o *Options) {
		o.InputType = inputType
	}
}

// WithOutputType creates outputType option
func WithOutputType(outputType reflect.Type) Option {
	return func(o *Options) {
		o.OutputType = outputType
	}
}

// WithLookupType creates lookupType option
func WithLookupType(fn func(name string) (reflect.Type, error)) Option {
	return func(o *Options) {
		o.LookupType = fn
	}
}

// NewOptions creates an options
func NewOptions(opts []Option) *Options {
	options := &Options{}
	options.apply(opts)
	return options
}

func (o *Options) apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}
