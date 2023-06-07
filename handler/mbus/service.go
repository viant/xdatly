package mbus

type Service interface {
	Push(message *Message) (*Confirmation, error)
}
