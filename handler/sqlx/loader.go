package sqlx

import "context"

type Loader interface {
	Flusher
	Load(ctx context.Context, tableName string, data interface{}) error
}
