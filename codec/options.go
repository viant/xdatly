package codec

import (
	"context"
	"io/fs"
	"reflect"
)

// Options is the reduced public codec option bag.
//
// It exists to support codec construction and execution only. It must not turn
// into a general-purpose runtime dependency transport for unrelated services.
type Options struct {
	LookupType  func(name string) (reflect.Type, error)
	Record      interface{}
	LookupValue func(ctx context.Context, name string) (interface{}, error)
	ResourceFS  fs.FS
	ColumnsSource
	Selector
	ValueGetter
	Options []interface{}
}

// WithResourceFS supplies package resources to codecs that compile external
// assets such as StructQL. The original fs.FS, including embed.FS, is retained.
func WithResourceFS(resources fs.FS) Option {
	return func(o *Options) {
		o.ResourceFS = resources
	}
}

// Apply applies option functions to the option bag.
func (o *Options) Apply(opts []Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// NewOptions creates a new option bag from option functions.
func NewOptions(opts []Option) *Options {
	options := &Options{}
	options.Apply(opts)
	return options
}

// Option mutates codec options.
type Option func(o *Options)

// ColumnsSource is the reduced public column metadata ask surface.
type ColumnsSource interface {
	Column(key string) (Column, bool)
	ColumnName(key string) (string, error)
}

// Column is the reduced public column metadata contract.
type Column interface {
	ColumnName() string
	ColumnType() reflect.Type
	FieldName() string
}

// ValueGetter is the reduced public value lookup contract.
type ValueGetter interface {
	Value(ctx context.Context, paramName string) (interface{}, error)
}

// Selector is the reduced public selector ask contract.
type Selector interface {
	IgnoreRead()
}

func WithTypeLookup(fn func(name string) (reflect.Type, error)) Option {
	return func(o *Options) {
		o.LookupType = fn
	}
}

func WithColumnsSource(columnSource ColumnsSource) Option {
	return func(o *Options) {
		o.ColumnsSource = columnSource
	}
}

func WithSelector(selector Selector) Option {
	return func(o *Options) {
		o.Selector = selector
	}
}

func WithRecord(record interface{}) Option {
	return func(o *Options) {
		o.Record = record
	}
}

func WithValueLookup(fn func(ctx context.Context, name string) (interface{}, error)) Option {
	return func(o *Options) {
		o.LookupValue = fn
	}
}

func WithValueGetter(valueGetter ValueGetter) Option {
	return func(o *Options) {
		o.ValueGetter = valueGetter
	}
}

func WithOptions(options ...interface{}) Option {
	return func(o *Options) {
		o.Options = options
	}
}
