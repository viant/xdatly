package state

import "context"

type Stater interface {
	Into(ctx context.Context, state interface{}, opt ...Option) error

	Value(ctx context.Context, key string) (interface{}, bool, error)
}

type Service struct {
	stater Stater
}

func (s Service) Value(ctx context.Context, key string) (interface{}, bool, error) {
	return s.stater.Value(ctx, key)
}

func (s Service) Into(ctx context.Context, state interface{}, opt ...Option) error {
	return s.stater.Into(ctx, state, opt...)
}

func New(stater Stater) *Service {
	return &Service{stater: stater}
}
