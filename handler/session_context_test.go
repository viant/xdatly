package handler

import (
	"context"
	"testing"

	"github.com/viant/xdatly/response"
)

type testContextBinder struct{}

func (b *testContextBinder) Bind(_ context.Context, target any) error {
	_ = target
	return nil
}

func (b *testContextBinder) Lookup(_ context.Context, key ValueKey) (any, bool, error) {
	_ = key
	return nil, false, nil
}

type testContextResponse struct{}

func (r *testContextResponse) StatusCode() int                   { return 0 }
func (r *testContextResponse) SetStatusCode(code int)            { _ = code }
func (r *testContextResponse) AddError(err error)                { _ = err }
func (r *testContextResponse) AddMetric(metric *response.Metric) { _ = metric }
func (r *testContextResponse) Metrics() response.Metrics         { return nil }

type testContextSession struct {
	binder   Binder
	response response.Writer
}

func (s *testContextSession) Binder() Binder            { return s.binder }
func (s *testContextSession) Response() response.Writer { return s.response }

func TestWithSessionAndSessionFromContext(t *testing.T) {
	sess := &testContextSession{
		binder:   &testContextBinder{},
		response: &testContextResponse{},
	}
	ctx := WithSession(context.Background(), sess)
	actual, ok := SessionFromContext(ctx)
	if !ok {
		t.Fatalf("expected session in context")
	}
	if actual != sess {
		t.Fatalf("expected stored session to round-trip")
	}
}

func TestWithSessionNilContext(t *testing.T) {
	sess := &testContextSession{
		binder:   &testContextBinder{},
		response: &testContextResponse{},
	}
	ctx := WithSession(nil, sess)
	if ctx == nil {
		t.Fatalf("expected non-nil context")
	}
	actual, ok := SessionFromContext(ctx)
	if !ok || actual != sess {
		t.Fatalf("expected session to be retrievable from generated background context")
	}
}

func TestSessionFromContextMissing(t *testing.T) {
	if actual, ok := SessionFromContext(context.Background()); ok || actual != nil {
		t.Fatalf("expected missing session")
	}
}

func TestWithSessionNilSession(t *testing.T) {
	ctx := context.Background()
	if actual := WithSession(ctx, nil); actual != ctx {
		t.Fatalf("expected nil session to leave context unchanged")
	}
}
