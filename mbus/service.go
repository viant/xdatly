package mbus

import "context"

// Message is the reduced public message-bus payload contract.
type Message struct {
	ID         string
	Resource   string
	TraceID    string
	Attributes map[string]interface{}
	Subject    string
	Data       interface{}
}

func (m *Message) AddAttribute(name string, value interface{}) {
	if len(m.Attributes) == 0 {
		m.Attributes = make(map[string]interface{})
	}
	m.Attributes[name] = value
}

// Service is the reduced public message-bus contract.
type Service interface {
	Push(ctx context.Context, message *Message) (*Confirmation, error)
	Message(dest string, data interface{}, opts ...Option) *Message
}

// Confirmation is the reduced delivery confirmation payload.
type Confirmation struct {
	MessageID string
}

func (c Confirmation) String() string {
	return c.MessageID
}

// Option mutates a message before publish.
type Option func(m *Message)

// Options applies a message option list.
type Options []Option

func (o Options) Apply(message *Message) {
	if len(o) == 0 {
		return
	}
	for _, opt := range o {
		opt(message)
	}
}

func WithAttribute(name string, value interface{}) Option {
	return func(m *Message) {
		m.AddAttribute(name, value)
	}
}

func WithTraceID(traceID string) Option {
	return func(m *Message) {
		m.TraceID = traceID
	}
}

func WithID(id string) Option {
	return func(m *Message) {
		m.ID = id
	}
}

func WithSubject(subject string) Option {
	return func(m *Message) {
		m.Subject = subject
	}
}
