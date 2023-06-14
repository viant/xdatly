package mbus

import "context"

type Mbus interface {
	Push(ctx context.Context, message *Message) (*Confirmation, error)
	Message(dest string, data interface{}, opts ...Option) *Message
}
type Service struct {
	mbus Mbus
}

func (s *Service) Push(ctx context.Context, message *Message) (*Confirmation, error) {
	return s.mbus.Push(ctx, message)
}

func (s *Service) Message(dest string, data interface{}, opts ...Option) *Message {
	return s.mbus.Message(dest, data, opts...)
}

func New(mbus Mbus) *Service {
	return &Service{mbus: mbus}
}
