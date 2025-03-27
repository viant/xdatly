package state

import (
	"context"
)

type Stater interface {

	//deprecated use bind instead
	Into(ctx context.Context, any interface{}, opt ...Option) error

	Bind(ctx context.Context, any interface{}, opt ...Option) error

	Value(ctx context.Context, key string) (interface{}, bool, error)

	ValuesOf(ctx context.Context, any interface{}) (map[string]interface{}, error)
}

type Service struct {
	stater Stater
}

func (s Service) Value(ctx context.Context, key string) (interface{}, bool, error) {
	return s.stater.Value(ctx, key)
}

// deprecated use bind instead
func (s Service) Into(ctx context.Context, state interface{}, opt ...Option) error {
	return s.stater.Into(ctx, state, opt...)
}

// Bind binds state based on parameter annotation
func (s Service) Bind(ctx context.Context, state interface{}, opt ...Option) error {
	return s.stater.Bind(ctx, state, opt...)
}

func New(stater Stater) *Service {
	return &Service{stater: stater}
}
