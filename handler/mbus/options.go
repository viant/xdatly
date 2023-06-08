package mbus

type Option func(m *Message)

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

func WithTraceId(traceID string) Option {
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
