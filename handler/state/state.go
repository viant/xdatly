package state

import "context"

type Stater interface {
	Into(ctx context.Context, state interface{}) error
}

type Service struct {
	stater Stater
}

func (s Service) Into(ctx context.Context, state interface{}) error {
	return s.stater.Into(ctx, state)
}

func New(stater Stater) *Service {
	return &Service{stater: stater}
}
