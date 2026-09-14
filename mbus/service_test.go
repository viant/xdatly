package mbus

import (
	"context"
	"testing"
)

type testService struct{}

func (testService) Push(_ context.Context, message *Message) (*Confirmation, error) {
	return &Confirmation{MessageID: message.ID}, nil
}

func (testService) Message(dest string, data interface{}, opts ...Option) *Message {
	msg := &Message{Resource: dest, Data: data}
	Options(opts).Apply(msg)
	return msg
}

func TestMessageOptionsAndServiceContracts(t *testing.T) {
	var _ Service = testService{}

	svc := testService{}
	msg := svc.Message("queue://demo", "payload", WithID("m1"), WithTraceID("t1"), WithSubject("sub"), WithAttribute("k", "v"))
	if msg.ID != "m1" || msg.TraceID != "t1" || msg.Subject != "sub" || msg.Resource != "queue://demo" {
		t.Fatalf("expected message fields to be preserved")
	}
	if msg.Attributes["k"] != "v" {
		t.Fatalf("expected message attributes to be preserved")
	}

	confirmation, err := svc.Push(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected push error: %v", err)
	}
	if confirmation.MessageID != "m1" {
		t.Fatalf("expected confirmation message id to match")
	}
}
