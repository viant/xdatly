package sqlx

import "context"

type Sequencer interface {
	Allocate(ctx context.Context, tableName string, dest interface{}, selector string) error
}
