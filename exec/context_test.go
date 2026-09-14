package exec

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/viant/xdatly/response"
	"github.com/viant/xdatly/tracing"
)

func TestGetContext(t *testing.T) {
	if GetContext(context.Background()) != nil {
		t.Fatalf("expected nil context lookup")
	}

	execCtx := NewContext("GET", "/v1/api/users", http.Header{"X-Trace": []string{"abc"}}, "v1")
	ctx := WithContext(context.Background(), execCtx)
	if GetContext(ctx) != execCtx {
		t.Fatalf("expected stored exec context")
	}
}

func TestNewGenericContext(t *testing.T) {
	trace := tracing.NewTrace("svc", "v1")
	start := time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)
	execCtx := New(
		WithMethod("JOB"),
		WithURI("job://sync"),
		WithTrace(trace),
		WithStartTime(start),
		WithHeaders(map[string]string{"X-Test": "1"}),
	)
	if execCtx.Method != "JOB" || execCtx.URI != "job://sync" {
		t.Fatalf("expected generic method/uri to be preserved, got %+v", execCtx)
	}
	if execCtx.Trace != trace || execCtx.TraceID != trace.TraceID {
		t.Fatalf("expected explicit trace to be preserved")
	}
	if !execCtx.StartTime.Equal(start) {
		t.Fatalf("expected explicit start time to be preserved")
	}
	if execCtx.Header["X-Test"] != "1" {
		t.Fatalf("expected generic headers to be copied")
	}
	if execCtx.Header == nil || execCtx.values == nil {
		t.Fatalf("expected generic context internals to initialize")
	}
}

func TestNewHTTPContextPreservesCompatibilityShape(t *testing.T) {
	execCtx := NewHTTPContext("GET", "/v1/api/users", http.Header{"X-Test": []string{"1"}}, "v1")
	if execCtx.Method != "GET" || execCtx.URI != "/v1/api/users" {
		t.Fatalf("expected http context to preserve method/uri")
	}
	if execCtx.Trace == nil || len(execCtx.Trace.Spans) == 0 {
		t.Fatalf("expected http context to append root span")
	}
}

func TestContextSnapshotForLogging(t *testing.T) {
	execCtx := NewContext("GET", "/v1/api/users", http.Header{"X-Test": []string{"1"}}, "v1")
	execCtx.AppendMetrics(&response.Metric{View: "users"})
	execCtx.SetError(errors.New("boom"))

	snapshot := execCtx.SnapshotForLogging()
	if snapshot == nil || len(snapshot.Metrics) != 1 || snapshot.Error != "boom" {
		t.Fatalf("expected snapshot to preserve metrics and error")
	}

	snapshot.Header["X-Test"] = "2"
	snapshot.Metrics[0].View = "changed"
	if execCtx.Header["X-Test"] != "1" {
		t.Fatalf("expected header deep copy")
	}
	if execCtx.Metrics[0].View != "users" {
		t.Fatalf("expected metrics deep copy")
	}
}

func TestPrepareLogging(t *testing.T) {
	execCtx := NewContext("GET", "/v1/api/users", http.Header{"X-Test": []string{"1"}}, "v1")
	execCtx.StatusCode = 200
	execCtx.AppendMetrics(&response.Metric{
		View:      "users",
		Type:      "SELECT",
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Executions: response.SQLExecutions{
			&response.SQLExecution{
				SQL:  "SELECT * FROM users WHERE id = ?",
				Args: []interface{}{1},
			},
		},
	})

	audit, trace := PrepareLogging(execCtx, false)
	if audit == nil || trace == nil {
		t.Fatalf("expected audit and trace snapshots")
	}
	if len(audit.Metrics) != 1 || audit.Metrics[0].Executions[0].SQL != "" {
		t.Fatalf("expected SQL to be hidden in audit metrics")
	}
	if len(trace.Spans) < 2 {
		t.Fatalf("expected root span plus metric-derived span")
	}
	if trace.Spans[0].Status.Code != "OK" {
		t.Fatalf("expected root status OK")
	}
}

func TestPrepareLogging_WithErrorSetsRootStatus(t *testing.T) {
	execCtx := NewContext("GET", "/v1/api/users", nil, "v1")
	execCtx.StatusCode = 500
	execCtx.SetError(errors.New("boom"))

	_, trace := PrepareLogging(execCtx, true)
	if trace == nil || len(trace.Spans) == 0 {
		t.Fatalf("expected trace snapshot")
	}
	if trace.Spans[0].Status.Code != "Error" || trace.Spans[0].Status.Message != "boom" {
		t.Fatalf("expected root error status")
	}
}

func TestNewContext_UsesIncomingTraceHeader(t *testing.T) {
	t.Setenv(trackingHeaderEnvKey, "X-Trace-Id")
	execCtx := NewContext("GET", "/v1/api/users", http.Header{
		"X-Trace-Id":    []string{"trace-123"},
		"Authorization": []string{"secret"},
	}, "v1")
	if execCtx.TraceID != "trace-123" {
		t.Fatalf("expected incoming trace id, got %q", execCtx.TraceID)
	}
	if _, ok := execCtx.Header["Authorization"]; ok {
		t.Fatalf("expected auth header to be omitted")
	}
}

func TestNewContext_UsesIncomingTraceparentHeader(t *testing.T) {
	execCtx := NewContext("GET", "/v1/api/users", http.Header{
		"Traceparent": []string{"00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"},
	}, "v1")
	if execCtx.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Fatalf("expected traceparent trace id, got %q", execCtx.TraceID)
	}
	if execCtx.Trace == nil || len(execCtx.Trace.Spans) == 0 {
		t.Fatalf("expected root trace span")
	}
	root := execCtx.Trace.Spans[0]
	if root.ParentSpanID == nil || *root.ParentSpanID != "00f067aa0ba902b7" {
		t.Fatalf("expected root parent span id from traceparent, got %+v", root.ParentSpanID)
	}
}
