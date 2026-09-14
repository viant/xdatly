package predicate

import "context"

// Criteria is the SQL fragment produced by a custom predicate handler.
type Criteria struct {
	Expression   string
	Placeholders []any
}

// Handler computes one predicate fragment from a bound input value.
type Handler interface {
	Compute(context.Context, any) (*Criteria, error)
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(context.Context, any) (*Criteria, error)

func (f HandlerFunc) Compute(ctx context.Context, value any) (*Criteria, error) {
	return f(ctx, value)
}
