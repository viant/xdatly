package differ

import "context"

type Service interface {
	Diff(ctx context.Context, from, to interface{}, options ...Option) *ChangeLog
}
