package http

import "context"

// Provider resolves or reuses a client from deployment-owned typed options.
// Runtime composition injects implementations. Callers do not own the provider
// or client and must not close their connection pools.
type Provider interface {
	Client(ctx context.Context, options Options) (Client, error)
}
