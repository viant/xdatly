package state

import (
	"context"
)

type Injector interface {

	//deprecated use bind instead
	Into(ctx context.Context, any interface{}, opt ...Option) error

	Bind(ctx context.Context, any interface{}, opt ...Option) error

	Value(ctx context.Context, key string) (interface{}, bool, error)

	ValuesOf(ctx context.Context, any interface{}) (map[string]interface{}, error)
}

type Service struct {
	injector Injector
}

func (s Service) Value(ctx context.Context, key string) (interface{}, bool, error) {
	return s.injector.Value(ctx, key)
}

// deprecated use bind instead
func (s Service) Into(ctx context.Context, state interface{}, opt ...Option) error {
	return s.injector.Into(ctx, state, opt...)
}

// Bind binds state based on parameter annotation
func (s Service) Bind(ctx context.Context, state interface{}, opt ...Option) error {
	return s.injector.Bind(ctx, state, opt...)
}

func (s Service) ValuesOf(ctx context.Context, any interface{}) (map[string]interface{}, error) {
	return s.injector.ValuesOf(ctx, any)
}

func New(stater Injector) *Service {
	return &Service{injector: stater}
}
