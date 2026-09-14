package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/viant/xdatly/response"
)

type testBinder struct {
	lastTarget any
	values     map[ValueKey]any
}

func (b *testBinder) Bind(_ context.Context, target any) error {
	b.lastTarget = target
	return nil
}

func (b *testBinder) Lookup(_ context.Context, key ValueKey) (any, bool, error) {
	value, ok := b.values[key]
	return value, ok, nil
}

type testResponse struct {
	status  int
	errors  []error
	metrics response.Metrics
}

func (r *testResponse) StatusCode() int {
	return r.status
}

func (r *testResponse) SetStatusCode(code int) {
	r.status = code
}

func (r *testResponse) AddError(err error) {
	r.errors = append(r.errors, err)
}

func (r *testResponse) AddMetric(metric *response.Metric) {
	r.metrics.Append(metric)
}

func (r *testResponse) Metrics() response.Metrics {
	return r.metrics
}

type testSession struct {
	binder   Binder
	response response.Writer
}

func (s *testSession) Binder() Binder {
	return s.binder
}

func (s *testSession) Response() response.Writer {
	return s.response
}

func TestSessionFuncExec(t *testing.T) {
	type input struct {
		Name string
	}
	type output struct {
		Seen string
	}

	binder := &testBinder{values: map[ValueKey]any{InputKey: &input{Name: "x"}}}
	resp := &testResponse{}
	sess := &testSession{binder: binder, response: resp}

	handler := ContractFunc[input, output](func(ctx context.Context, sess Session, in *input, out *output) error {
		if _, ok, err := sess.Binder().Lookup(ctx, InputKey); err != nil || !ok {
			t.Fatalf("expected InputKey lookup to succeed, err=%v ok=%v", err, ok)
		}
		out.Seen = in.Name
		sess.Response().SetStatusCode(202)
		return nil
	})

	out := &output{}
	if err := handler.Exec(context.Background(), sess, &input{Name: "bound"}, out); err != nil {
		t.Fatalf("unexpected exec error: %v", err)
	}
	if out.Seen != "bound" {
		t.Fatalf("expected output to be updated, got %q", out.Seen)
	}
	if resp.StatusCode() != 202 {
		t.Fatalf("expected status 202, got %d", resp.StatusCode())
	}
}

func TestBinderCompatibilityShape(t *testing.T) {
	type generatedInput struct {
		Name string
	}
	initLike := func(ctx context.Context, sess Session, input *generatedInput) error {
		return sess.Binder().Bind(ctx, input)
	}

	binder := &testBinder{}
	sess := &testSession{
		binder:   binder,
		response: &testResponse{},
	}
	input := &generatedInput{Name: "demo"}
	if err := initLike(context.Background(), sess, input); err != nil {
		t.Fatalf("unexpected bind error: %v", err)
	}
	if binder.lastTarget != input {
		t.Fatalf("expected binder to receive the generated input target")
	}
}

func TestBinderLookupUnknownKeyShape(t *testing.T) {
	binder := &testBinder{values: map[ValueKey]any{}}
	value, ok, err := binder.Lookup(context.Background(), ValueKey("missing"))
	if err != nil {
		t.Fatalf("unexpected lookup error: %v", err)
	}
	if ok {
		t.Fatalf("expected unknown key lookup to report ok=false")
	}
	if value != nil {
		t.Fatalf("expected unknown key lookup to return nil value")
	}
}

func TestResponseWriterRecordsErrorAndStatus(t *testing.T) {
	resp := &testResponse{}
	errExpected := errors.New("boom")
	resp.SetStatusCode(409)
	resp.AddError(errExpected)
	resp.AddMetric(&response.Metric{View: "users"})

	if resp.StatusCode() != 409 {
		t.Fatalf("expected status 409, got %d", resp.StatusCode())
	}
	if len(resp.errors) != 1 {
		t.Fatalf("expected one recorded error, got %d", len(resp.errors))
	}
	if resp.errors[0] != errExpected {
		t.Fatalf("expected recorded error to match")
	}
	if len(resp.Metrics()) != 1 || resp.Metrics()[0].View != "users" {
		t.Fatalf("expected one recorded metric")
	}
}
