package sqlx

import "context"

type Reader interface {
	Read(ctx context.Context, dest interface{}, SQL string, params ...interface{}) error
}

type ReadOptions struct {
	Args []interface{}
}

type ReaderOption func(o *ReadOptions)
