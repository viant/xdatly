package docs

import "context"

// Provider creates a documentation lookup service.
type Provider interface {
	Service(ctx context.Context, options ...Option) (Service, error)
}

// Service looks up documentation payloads by key.
type Service interface {
	Lookup(ctx context.Context, key string) (string, bool, error)
}
