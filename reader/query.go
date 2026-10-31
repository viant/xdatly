package reader

import (
	"context"
	"database/sql"
)

// ParameterizedQuery pairs a complete SQL statement with its ordered bound
// arguments. Query retains placeholder markers; Args contains their values.
// Callers must not interpolate the values into Query.
type ParameterizedQuery struct {
	Query string
	Args  []any
}

// QueryGuard optionally admits a fully bound query on its source connection.
// A nil error permits execution; an error prevents that statement. Implementations
// must be concurrency-safe and must not mutate the query or execute it themselves.
type QueryGuard func(context.Context, *sql.DB, *ParameterizedQuery) error
