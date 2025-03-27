package tracing

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
	StatusUnset = "Unset"
)

type Trace struct {
	TraceID  string        `json:"traceId"`
	Spans    []*Span       `json:"spans"`
	Resource *ResourceInfo `json:"resource"`
}

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

func (s *Span) OnDone() {
	s.EndTime = time.Now()
}

type SpanStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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
		Attributes:   make(map[string]string),
		Status: SpanStatus{
			Code:    StatusUnset,
			Message: "",
		},
	}
}

// SetStatus updates the status of the span based on the provided error.
func (s *Span) SetStatus(err error) {
	if err != nil {
		s.Status = SpanStatus{
			Code:    StatusError,
			Message: fmt.Sprintf(err.Error()),
		}
	} else {
		s.Status = SpanStatus{
			Code:    StatusOK,
			Message: "",
		}
	}
}

// SetStatusFromHTTPCode updates the span's status based on the HTTP response code.
func (s *Span) SetStatusFromHTTPCode(code int) {
	switch {
	case code >= 100 && code < 400:
		// For HTTP status codes in the 1xx, 2xx, or 3xx range, status is unset.
		s.Status = SpanStatus{
			Code:    StatusOK, // Unset
			Message: "",
		}
	case code >= 400 && code < 500:
		// For HTTP status codes in the 4xx range:
		if s.Kind == "CLIENT" {
			// Client spans should be set to Error.
			s.Status = SpanStatus{
				Code:    StatusError, // Error
				Message: http.StatusText(code),
			}
		} else {
			// Server spans should remain unset.
			s.Status = SpanStatus{
				Code:    StatusOK, // Unset
				Message: "",
			}
		}
	case code >= 500:
		// For HTTP status codes in the 5xx range, status is set to Error.
		s.Status = SpanStatus{
			Code:    StatusError, // Error
			Message: http.StatusText(code),
		}
	default:
		// For any other codes, default to unset status.
		s.Status = SpanStatus{
			Code:    StatusOK, // Unset
			Message: "",
		}
	}
}

func (t *Trace) Append(span ...*Span) {
	t.Spans = append(t.Spans, span...)
}

func (s *Span) WithAttributes(attrs map[string]string) *Span {
	for k, v := range attrs {
		s.Attributes[k] = v
	}
	return s
}
