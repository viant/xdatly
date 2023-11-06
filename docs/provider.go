package docs

import "context"

type (
	Provider interface {
		Service(ctx context.Context, options ...Option) (Service, error)
	}

	Service interface {
		Lookup(ctx context.Context, key string) (string, bool, error)
	}
)
