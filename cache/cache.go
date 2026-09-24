// Package cache defines an optional, backend-neutral byte cache capability.
// Callers own key construction, scope, and authorization policy. Cache and
// Provider values are borrowed from composition; callers must not close them.
package cache

import (
	"context"
	"time"
)

// Cache stores opaque bytes. A Get miss returns found=false and a nil error.
// Implementations must not expose mutable storage through Put or Get. A
// positive TTL expires an entry; zero means no expiry.
type Cache interface {
	Get(ctx context.Context, key string) (value []byte, found bool, err error)
	Put(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// Provider resolves exactly the named, explicitly registered cache. Empty or
// unknown names fail; resolution does not create a backend.
type Provider interface {
	Cache(ctx context.Context, name string) (Cache, error)
}
