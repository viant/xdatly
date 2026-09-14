package exec

import (
	"net/http"

	"github.com/viant/xdatly/tracing"
)

// NewHTTPContext builds one execution context from HTTP request metadata.
// This helper is transport-specific; the stable generic entrypoint is New(...).
func NewHTTPContext(method string, URI string, header http.Header, version string) *Context {
	ret := New(
		WithMethod(method),
		WithURI(URI),
		WithTraceResource("datly", version),
	)
	ret.setHeader(header)
	if ret.Trace == nil {
		ret.Trace = tracing.NewTrace("datly", version)
	}
	root := tracing.NewSpan("HTTP "+method+" "+URI, "SERVER", ret.parentSpanID, ret.StartTime, ret.StartTime)
	root.WithAttributes(map[string]string{
		"http.method": method,
		"http.url":    URI,
	})
	ret.Trace.Append(&root)
	if ret.TraceID != "" {
		ret.Trace.TraceID = ret.TraceID
	} else {
		ret.TraceID = ret.Trace.TraceID
	}
	return ret
}

// NewContext is the compatibility wrapper over the older HTTP-colored
// constructor signature. Prefer New(...) for generic callers and
// NewHTTPContext(...) for explicit HTTP assembly.
func NewContext(method string, URI string, header http.Header, version string) *Context {
	return NewHTTPContext(method, URI, header, version)
}
