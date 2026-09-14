package codec

import (
	"context"
	"reflect"
)

// Config describes the exact source and destination contract for a codec.
type Config struct {
	Body                 string
	SourceType           reflect.Type `json:"-" yaml:"-"`
	DestinationType      reflect.Type `json:"-" yaml:"-"`
	Args                 []string
	OutputTypeExpression string
}

// Factory creates codec instances from config plus options during registration.
type Factory interface {
	New(codecConfig *Config, options ...Option) (Instance, error)
}

// Instance is the reduced public codec execution contract. Registered instances
// may be shared by concurrent invocations: Value must be concurrency-safe and
// must not retain invocation context, record pointers, or lookup options.
type Instance interface {
	Value(ctx context.Context, raw interface{}, options ...Option) (interface{}, error)
}
