// Package connector defines an explicitly opt-in escape hatch for specialized
// handler integrations. It is separate from managed Data and transaction scopes.
package connector

import (
	"context"
	"database/sql"
)

// Provider resolves an exact registered name. Empty/unknown names and canceled
// contexts must fail; implementations must never silently select a default DB.
// The returned DB is borrowed: handlers must not close it. Direct work through
// it is outside Datly's managed Data transaction and flush lifecycle.
type Provider interface {
	Connector(ctx context.Context, name string) (*sql.DB, error)
}
