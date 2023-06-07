package sqlx

import "context"

type Flusher interface {
	Flush(ctx context.Context, tableName string) error
}
