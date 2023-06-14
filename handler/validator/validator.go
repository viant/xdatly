package validator

import "context"

type Validator interface {
	Validate(ctx context.Context, any interface{}, options ...Option) *Validation
}

type Service struct {
	validator Validator
}

func (s *Service) Validate(ctx context.Context, any interface{}, options ...Option) *Validation {
	return s.validator.Validate(ctx, any, options...)
}

func New(validator Validator) *Service {
	return &Service{validator: validator}
}
