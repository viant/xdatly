package exec

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/viant/scy/auth/jwt"
	"github.com/viant/xdatly/handler/async"
	"github.com/viant/xdatly/handler/response"
	"github.com/viant/xdatly/handler/tracing"
)

type contextKey string
type errorKey string

var ContextKey = contextKey("context")
var ErrorKey = errorKey("error")

func GetContext(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(ContextKey)
	if value == nil {
		return nil
	}
	return value.(*Context)
}

// Context represents an execution context
type Context struct {
	Method                     string            `json:"method,omitempty"`
	URI                        string            `json:"uri,omitempty"`
	StatusCode                 int               `json:"statusCode,omitempty"`
	Status                     string            `json:"status,omitempty"`
	Error                      string            `json:"error,omitempty"`
	ElapsedMs                  int               `json:"elapsedMs,omitempty"`
	StartTime                  time.Time         `json:"startTime,omitempty"`
	Auth                       *jwt.Claims       `json:"auth,omitempty"`
	Header                     map[string]string `json:"header,omitempty"`
	Metrics                    response.Metrics  `json:"metrics,omitempty"`
	TraceID                    string            `json:"traceId,omitempty"`
	Trace                      *tracing.Trace    `json:"-"`
	mux                        sync.RWMutex
	jobs                       []*async.Job
	values                     map[string]interface{}
	IgnoreEmptyQueryParameters bool `json:"-"`
}

func (c *Context) AppendMetrics(metrics *response.Metric) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Metrics = append(c.Metrics, metrics)
}

func (c *Context) SetError(err error) {
	c.Error = err.Error()
	c.Status = "error"
}

const trackingHeaderEnvKey = "XDATLY_TRACING_HEADER"

func (c *Context) setHeader(header http.Header) {
	c.Header = make(map[string]string)
	trackingHeaderKey := os.Getenv(trackingHeaderEnvKey)
	if trackingHeaderKey == "" {
		trackingHeaderKey = "xtraceid"
	}
	trackingHeaderKey = strings.ReplaceAll(strings.ToLower(trackingHeaderKey), "-", "")
	for k := range header {
		lowerKey := strings.ToLower(k)
		if strings.Contains(lowerKey, "auth") {
			continue
		}
		if trackingHeaderKey == strings.ReplaceAll(lowerKey, "-", "") {
			c.TraceID = header.Get(k)
			continue
		}
		c.Header[k] = header.Get(k)
	}
}

func (c *Context) SetValue(key string, value interface{}) {
	c.mux.Lock()
	c.values[key] = value
	c.mux.Unlock()
}

func (c *Context) Value(key string) (interface{}, bool) {
	c.mux.RLock()
	value, has := c.values[key]
	c.mux.RUnlock()
	return value, has
}

func (c *Context) Elapsed() time.Duration {
	now := time.Now()
	return now.Sub(c.StartTime)
}

func (c *Context) EndTime() time.Time {
	now := time.Now()
	return now
}

func (c *Context) AsyncElapsed() time.Duration {
	if len(c.jobs) == 0 {
		return 0
	}
	started := c.jobs[0].CreationTime
	ended := started

	for _, job := range c.jobs {
		if job.CreationTime.Before(started) {
			started = job.CreationTime
		}
		if job.EndTime != nil && job.EndTime.After(ended) {
			ended = *job.EndTime
		}
	}
	if ended == started {
		ended = time.Now()
	}
	return ended.Sub(started)
}

func (c *Context) AsyncEndTime() *time.Time {
	if len(c.jobs) == 0 {
		return nil
	}
	var ret *time.Time
	for _, job := range c.jobs {
		if job.EndTime != nil {
			if ret == nil {
				ret = job.EndTime
			} else if job.EndTime.After(*ret) {
				ret = job.EndTime
			}
		}
	}
	return ret
}

func (c *Context) AsyncCreationTime() *time.Time {
	if len(c.jobs) == 0 {
		return nil
	}
	var ret *time.Time
	for _, job := range c.jobs {
		if ret == nil {
			ret = &job.CreationTime
		} else if job.CreationTime.Before(*ret) {
			ret = &job.CreationTime
		}

	}
	return ret
}

func (c *Context) AppendJob(job *async.Job) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.hasJob(job) {
		return
	}

	c.jobs = append(c.jobs, job)
}

func (c *Context) hasJob(job *async.Job) bool {
	for _, candidate := range c.jobs {
		if candidate.ID == job.ID {
			return true
		}
	}
	return false
}

func (c *Context) AsyncStatus() string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if len(c.jobs) == 0 {
		return "N/A"
	}
	pendingCount := 0
	runningCount := 0
	doneCount := 0
	for _, candidate := range c.jobs {
		if candidate.Status == string(async.StatusDone) || candidate.Status == string(async.StatusError) {
			doneCount++
		} else if candidate.Status == string(async.StatusRunning) {
			runningCount++
		} else if candidate.Status == string(async.StatusPending) {
			pendingCount++
		}
	}
	if doneCount == len(c.jobs) {
		return string(async.StatusDone)
	}
	if pendingCount == len(c.jobs) {
		return string(async.StatusPending)
	}
	return string(async.StatusRunning)
}

// CreateInitialSpan constructs the first span for the trace.
func CreateInitialSpan(method, uri string) *tracing.Span {
	startTime := time.Now()
	return &tracing.Span{
		SpanID:    uuid.New().String(),
		Name:      "HTTP " + method + " " + uri,
		Kind:      "SERVER",
		StartTime: startTime,
		EndTime:   startTime, // Update this when the operation completes
		Attributes: map[string]string{
			"http.method": method,
			"http.url":    uri,
			// Add other relevant attributes as needed
		},
		Status: tracing.SpanStatus{
			Code:    tracing.StatusOK,
			Message: "",
		},
	}
}

// NewContext creates a new context
func NewContext(method string, URI string, header http.Header, version string) *Context {
	trace := tracing.NewTrace("datly", version)
	ret := &Context{Method: method,
		URI:       URI,
		Trace:     trace,
		StartTime: time.Now(),
		values:    map[string]interface{}{}}

	ret.setHeader(header)
	trace.Append(CreateInitialSpan(method, URI))
	if ret.TraceID != "" {
		trace.TraceID = ret.TraceID
	} else {
		ret.TraceID = trace.TraceID
	}
	return ret
}
