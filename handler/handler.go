package handler

import "context"

type Handler interface {
	Exec(ctx context.Context, session Session) (interface{}, error)
}
