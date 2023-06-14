package differ

import "context"

type Differ interface {
	Diff(ctx context.Context, from, to interface{}, options ...Option) *ChangeLog
}

type Service struct {
	differ Differ
}

func (s *Service) Diff(ctx context.Context, from, to interface{}, options ...Option) *ChangeLog {
	return s.differ.Diff(ctx, from, to, options...)
}

func New(differ Differ) *Service {
	return &Service{differ: differ}
}
