package handler

import (
	"context"
	"testing"

	"github.com/viant/xdatly/response"
)

type hookTestBinder struct{}

func (b *hookTestBinder) Bind(_ context.Context, target any) error {
	_ = target
	return nil
}

func (b *hookTestBinder) Lookup(_ context.Context, key ValueKey) (any, bool, error) {
	_ = key
	return nil, false, nil
}

type hookTestResponse struct{}

func (r *hookTestResponse) StatusCode() int                   { return 0 }
func (r *hookTestResponse) SetStatusCode(code int)            { _ = code }
func (r *hookTestResponse) AddError(err error)                { _ = err }
func (r *hookTestResponse) AddMetric(metric *response.Metric) { _ = metric }
func (r *hookTestResponse) Metrics() response.Metrics         { return nil }

type hookTestSession struct {
	binder   Binder
	response response.Writer
}

func (s *hookTestSession) Binder() Binder            { return s.binder }
func (s *hookTestSession) Response() response.Writer { return s.response }

func TestWithHookAppliedAndHookApplied(t *testing.T) {
	if HookApplied(context.Background()) {
		t.Fatalf("expected plain context to have no hook marker")
	}
	ctx := WithHookApplied(context.Background())
	if !HookApplied(ctx) {
		t.Fatalf("expected hook marker to be present")
	}
}

func TestRunPreInvokeHook(t *testing.T) {
	var calls int
	err := RunPreInvokeHook(context.Background(), "input", func(ctx context.Context, input any) error {
		_ = ctx
		_ = input
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected hook error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected hook to run once, got %d", calls)
	}
}

func TestRunPreInvokeHookSkippedWhenAlreadyApplied(t *testing.T) {
	var calls int
	err := RunPreInvokeHook(WithHookApplied(context.Background()), "input", func(ctx context.Context, input any) error {
		_ = ctx
		_ = input
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected hook error: %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected hook to be skipped, got %d calls", calls)
	}
}

func TestInvokeWithSession(t *testing.T) {
	sess := &hookTestSession{
		binder:   &hookTestBinder{},
		response: &hookTestResponse{},
	}
	actual, err := InvokeWithSession(context.Background(), sess, "input", func(ctx context.Context, input any) (any, error) {
		if !HookApplied(ctx) {
			t.Fatalf("expected invoke context to carry hook-applied marker")
		}
		resolved, ok := SessionFromContext(ctx)
		if !ok || resolved != sess {
			t.Fatalf("expected invoke context to carry session")
		}
		return input, nil
	})
	if err != nil {
		t.Fatalf("unexpected invoke error: %v", err)
	}
	if actual != "input" {
		t.Fatalf("unexpected invoke result %v", actual)
	}
}
