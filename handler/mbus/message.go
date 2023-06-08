package mbus

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
