package validator

import "context"

type Service interface {
	Validate(ctx context.Context, any interface{}, options ...Option) *Validation
}
