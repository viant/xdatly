package handler

import "context"

// Contract is the minimal public custom-handler contract.
// Metadata remains the canonical source of handler identity; this is
// conformance-facing only.
type Contract[I any, O any] interface {
	Exec(ctx context.Context, sess Session, input *I, output *O) error
}

// ContractFunc adapts a function to Contract.
type ContractFunc[I any, O any] func(ctx context.Context, sess Session, input *I, output *O) error

func (f ContractFunc[I, O]) Exec(ctx context.Context, sess Session, input *I, output *O) error {
	return f(ctx, sess, input, output)
}
