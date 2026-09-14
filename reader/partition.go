// Package reader defines public extension contracts for typed Datly reads.
package reader

import (
	"context"
	"database/sql"
)

type PartitionRequest struct {
	DB        *sql.DB
	View      string
	Arguments []string
}

type Partition struct {
	Table      string
	Expression string
	Args       []any
}

type Partitioner interface {
	Partitions(ctx context.Context, request PartitionRequest) ([]Partition, error)
}

type ReducerProvider interface {
	Reducer(ctx context.Context) Reducer
}

type Reducer interface {
	// Reduce receives the complete typed result of a view fetch. Input order
	// is unspecified across partitions/batches; order-sensitive reducers must
	// sort explicitly. Relation readers invoke it once after all batches finish.
	Reduce(ctx context.Context, rows any) (any, error)
}
