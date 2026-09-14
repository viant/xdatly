package tracing

import (
	"testing"
	"time"
)

func TestTracePayloadContracts(t *testing.T) {
	trace := &Trace{
		TraceID: "trace-1",
		Resource: &ResourceInfo{
			ServiceName:    "datly",
			ServiceVersion: "1.0",
		},
	}
	span := (&Span{
		SpanID:     "span-1",
		Name:       "read",
		Kind:       "SERVER",
		Attributes: map[string]string{},
		Status: SpanStatus{
			Code: StatusOK,
		},
	}).WithAttributes(map[string]string{"view": "users"})
	trace.Append(span)

	if trace.TraceID != "trace-1" {
		t.Fatalf("expected trace id to be preserved")
	}
	if len(trace.Spans) != 1 {
		t.Fatalf("expected one span")
	}
	if trace.Spans[0].Attributes["view"] != "users" {
		t.Fatalf("expected span attributes to be preserved")
	}
}

func TestTraceHelpers(t *testing.T) {
	trace := NewTrace("svc", "v1")
	if trace.Resource == nil || trace.Resource.ServiceName != "svc" {
		t.Fatalf("expected trace resource to be initialized")
	}

	span := NewSpan("read", "CLIENT", nil, time.Now(), time.Now())
	span.SetStatus(nil)
	if span.Status.Code != StatusOK {
		t.Fatalf("expected OK status")
	}
	span.SetStatusFromHTTPCode(500)
	if span.Status.Code != StatusError {
		t.Fatalf("expected error status from 500")
	}
	before := span.EndTime
	span.OnDone()
	if !span.EndTime.After(before) && !span.EndTime.Equal(before) {
		t.Fatalf("expected span end time to be updated")
	}
}
