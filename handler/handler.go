package handler

import (
	"context"
)

type Handler interface {
	Exec(ctx context.Context, session Session) (interface{}, error)
}

type (
	Factory interface {
		New(ctx context.Context, options ...Option) (Handler, error)
	}
)
