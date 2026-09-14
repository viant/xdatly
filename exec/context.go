package exec

import (
	"context"
	"sync"
	"time"

	"github.com/viant/xdatly/response"
	"github.com/viant/xdatly/tracing"
)

type contextKey string

var ContextKey = contextKey("xdatly_exec_context")

func GetContext(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	ret, _ := value.(*Context)
	return ret
}

// WithContext attaches one execution context to ctx under the public ContextKey.
func WithContext(ctx context.Context, execCtx *Context) context.Context {
	if execCtx == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ContextKey, execCtx)
}

type Context struct {
	Method     string            `json:"method,omitempty"`
	URI        string            `json:"uri,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Status     string            `json:"status,omitempty"`
	Error      string            `json:"error,omitempty"`
	ElapsedMs  int               `json:"elapsedMs,omitempty"`
	StartTime  time.Time         `json:"startTime,omitempty"`
	Header     map[string]string `json:"header,omitempty"`
	Metrics    response.Metrics  `json:"metrics,omitempty"`
	TraceID    string            `json:"traceId,omitempty"`
	Trace      *tracing.Trace    `json:"-"`

	mux    sync.RWMutex
	values map[string]any

	parentSpanID *string
}

// Option mutates the generic execution context shape.
type Option func(*Context)

// New creates a generic execution context without assuming any specific
// transport. Callers add transport/runtime details through options.
func New(options ...Option) *Context {
	ret := &Context{
		values: map[string]any{},
	}
	for _, option := range options {
		if option != nil {
			option(ret)
		}
	}
	if ret.StartTime.IsZero() {
		ret.StartTime = time.Now()
	}
	if ret.Trace != nil && ret.TraceID == "" {
		ret.TraceID = ret.Trace.TraceID
	}
	if ret.Header == nil {
		ret.Header = map[string]string{}
	}
	return ret
}

func WithMethod(method string) Option {
	return func(c *Context) {
		c.Method = method
	}
}

func WithURI(uri string) Option {
	return func(c *Context) {
		c.URI = uri
	}
}

func WithStatusCode(code int) Option {
	return func(c *Context) {
		c.StatusCode = code
	}
}

func WithStatus(status string) Option {
	return func(c *Context) {
		c.Status = status
	}
}

func WithStartTime(start time.Time) Option {
	return func(c *Context) {
		c.StartTime = start
	}
}

func WithTrace(trace *tracing.Trace) Option {
	return func(c *Context) {
		c.Trace = trace
		if trace != nil {
			c.TraceID = trace.TraceID
		}
	}
}

func WithTraceResource(serviceName, serviceVersion string) Option {
	return func(c *Context) {
		trace := tracing.NewTrace(serviceName, serviceVersion)
		c.Trace = trace
		c.TraceID = trace.TraceID
	}
}

func WithHeaders(header map[string]string) Option {
	return func(c *Context) {
		if header == nil {
			c.Header = nil
			return
		}
		copyHeader := make(map[string]string, len(header))
		for k, v := range header {
			copyHeader[k] = v
		}
		c.Header = copyHeader
	}
}

func (c *Context) AppendMetrics(metric *response.Metric) {
	if c == nil || metric == nil {
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Metrics = append(c.Metrics, metric)
}

func (c *Context) SetError(err error) {
	if c == nil || err == nil {
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Error = err.Error()
	c.Status = "error"
}

func (c *Context) SetValue(key string, value any) {
	if c == nil {
		return
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.values == nil {
		c.values = map[string]any{}
	}
	c.values[key] = value
}

func (c *Context) Value(key string) (any, bool) {
	if c == nil {
		return nil, false
	}
	c.mux.RLock()
	defer c.mux.RUnlock()
	value, ok := c.values[key]
	return value, ok
}

func (c *Context) SnapshotForLogging() *Context {
	if c == nil {
		return nil
	}
	c.mux.RLock()
	defer c.mux.RUnlock()

	snapshot := &Context{
		Method:     c.Method,
		URI:        c.URI,
		StatusCode: c.StatusCode,
		Status:     c.Status,
		Error:      c.Error,
		ElapsedMs:  int(time.Since(c.StartTime).Milliseconds()),
		StartTime:  c.StartTime,
		TraceID:    c.TraceID,
	}

	if c.Header != nil {
		headerCopy := make(map[string]string, len(c.Header))
		for k, v := range c.Header {
			headerCopy[k] = v
		}
		snapshot.Header = headerCopy
	}
	if c.Metrics != nil {
		snapshot.Metrics = deepCopyMetrics(c.Metrics)
	}
	if c.Trace != nil {
		snapshot.Trace = deepCopyTrace(c.Trace)
	}
	if c.values != nil {
		valuesCopy := make(map[string]any, len(c.values))
		for k, v := range c.values {
			valuesCopy[k] = v
		}
		snapshot.values = valuesCopy
	}
	if c.parentSpanID != nil {
		parentSpanID := *c.parentSpanID
		snapshot.parentSpanID = &parentSpanID
	}
	return snapshot
}

func deepCopyMetrics(src response.Metrics) response.Metrics {
	if src == nil {
		return nil
	}
	result := make(response.Metrics, len(src))
	for i, metric := range src {
		if metric == nil {
			continue
		}
		clonedMetric := *metric
		if metric.Executions != nil {
			executionsCopy := make(response.SQLExecutions, len(metric.Executions))
			for j, exec := range metric.Executions {
				if exec == nil {
					continue
				}
				clonedExec := *exec
				if exec.Args != nil {
					argsCopy := make([]interface{}, len(exec.Args))
					copy(argsCopy, exec.Args)
					clonedExec.Args = argsCopy
				}
				if exec.CacheStats != nil {
					cacheCopy := *exec.CacheStats
					clonedExec.CacheStats = &cacheCopy
				}
				executionsCopy[j] = &clonedExec
			}
			clonedMetric.Executions = executionsCopy
		}
		result[i] = &clonedMetric
	}
	return result
}

func deepCopyTrace(src *tracing.Trace) *tracing.Trace {
	if src == nil {
		return nil
	}
	traceCopy := *src
	if src.Resource != nil {
		resourceCopy := *src.Resource
		traceCopy.Resource = &resourceCopy
	}
	if src.Spans != nil {
		spansCopy := make([]*tracing.Span, len(src.Spans))
		for i, span := range src.Spans {
			if span == nil {
				continue
			}
			spanCopy := *span
			if span.Attributes != nil {
				attrsCopy := make(map[string]string, len(span.Attributes))
				for k, v := range span.Attributes {
					attrsCopy[k] = v
				}
				spanCopy.Attributes = attrsCopy
			}
			spansCopy[i] = &spanCopy
		}
		traceCopy.Spans = spansCopy
	}
	return &traceCopy
}
