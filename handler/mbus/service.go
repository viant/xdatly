package mbus

import "context"

type Service interface {
	Push(ctx context.Context, message *Message) (*Confirmation, error)
	Message(dest string, data interface{}, opts ...Option) *Message
}
