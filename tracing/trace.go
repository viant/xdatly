package tracing

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
	StatusUnset = "Unset"
)

// Trace is the reduced public tracing payload contract.
type Trace struct {
	TraceID  string        `json:"traceId"`
	Spans    []*Span       `json:"spans"`
	Resource *ResourceInfo `json:"resource"`
}

// Span is the reduced public tracing span payload contract.
type Span struct {
	SpanID       string            `json:"spanId"`
	ParentSpanID *string           `json:"parentSpanId,omitempty"`
	Name         string            `json:"name"`
	Kind         string            `json:"kind"`
	StartTime    time.Time         `json:"startTime"`
	EndTime      time.Time         `json:"endTime"`
	Attributes   map[string]string `json:"attributes"`
	Status       SpanStatus        `json:"status"`
}

// SpanStatus is the reduced public span status payload.
type SpanStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResourceInfo is the reduced public trace resource payload.
type ResourceInfo struct {
	ServiceName    string `json:"service.name"`
	ServiceVersion string `json:"service.version"`
}

func NewTrace(serviceName, serviceVersion string) *Trace {
	return &Trace{
		TraceID: uuid.New().String(),
		Spans:   []*Span{},
		Resource: &ResourceInfo{
			ServiceName:    serviceName,
			ServiceVersion: serviceVersion,
		},
	}
}

func NewSpan(name, kind string, parentID *string, startTime, endTime time.Time) Span {
	return Span{
		SpanID:       uuid.New().String(),
		ParentSpanID: parentID,
		Name:         name,
		Kind:         kind,
		StartTime:    startTime,
		EndTime:      endTime,
		Attributes:   map[string]string{},
		Status: SpanStatus{
			Code:    StatusUnset,
			Message: "",
		},
	}
}

func (t *Trace) Append(span ...*Span) {
	t.Spans = append(t.Spans, span...)
}

func (s *Span) OnDone() {
	s.EndTime = time.Now()
}

func (s *Span) SetStatus(err error) {
	if err != nil {
		s.Status = SpanStatus{
			Code:    StatusError,
			Message: err.Error(),
		}
		return
	}
	s.Status = SpanStatus{
		Code:    StatusOK,
		Message: "",
	}
}

func (s *Span) WithAttributes(attrs map[string]string) *Span {
	if s.Attributes == nil {
		s.Attributes = map[string]string{}
	}
	for k, v := range attrs {
		s.Attributes[k] = v
	}
	return s
}
