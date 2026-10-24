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

// VariadicHandler receives trusted authored arguments in addition to the bound
// input. Implement this optional contract to accept handler configuration;
// existing Handler implementations and HandlerFunc remain unchanged.
type VariadicHandler interface {
	Compute(context.Context, any, ...any) (*Criteria, error)
}

// VariadicHandlerFunc adapts a function to VariadicHandler.
type VariadicHandlerFunc func(context.Context, any, ...any) (*Criteria, error)

func (f VariadicHandlerFunc) Compute(ctx context.Context, value any, args ...any) (*Criteria, error) {
	return f(ctx, value, args...)
}
