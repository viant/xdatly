package exec

import (
	"errors"

	"github.com/viant/xdatly/tracing"
)

// PrepareLogging is the compatibility helper that derives audit/trace-ready
// snapshots from one execution context. It is not the core execution-context
// contract itself; the stable core remains the Context payload plus its
// generic constructors and accessors.
func PrepareLogging(ctx *Context, includeSQL bool) (*Context, *tracing.Trace) {
	if ctx == nil {
		return nil, nil
	}
	snapshot := ctx.SnapshotForLogging()
	if snapshot == nil {
		return nil, nil
	}
	if !includeSQL {
		snapshot.Metrics = snapshot.Metrics.HideMetrics()
	}
	trace := snapshot.Trace
	if trace == nil || len(trace.Spans) == 0 {
		return snapshot, trace
	}
	rootSpan := trace.Spans[0]
	trace.Append(snapshot.Metrics.ToSpans(&rootSpan.SpanID)...)
	if snapshot.Error != "" {
		trace.Spans[0].SetStatus(errors.New(snapshot.Error))
	} else {
		trace.Spans[0].SetStatusFromHTTPCode(snapshot.StatusCode)
	}
	return snapshot, trace
}
