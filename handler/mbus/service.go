package mbus

type Service interface {
	Push(message *Message) (*Confirmation, error)
	Message(dest string, data interface{}, opts ...Option) *Message
}
